package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageLifecycle(t *testing.T) {
	user1 := "cutie"
	user2 := "boykisser"
	//create users
	cutieToken, err := loginAndGetToken(user1, "password123")
	require.NoError(t, err, "cutie failed to login")

	boykisserToken, err := loginAndGetToken(user2, "password123")
	require.NoError(t, err, "boykisser failed to login")

	//connect users to the websocket
	cutieWs, _, err := connectToWebSocket(cutieToken)
	require.NoError(t, err, "cutie failed to connect to the web socket")
	defer cutieWs.Close()

	boykisserWS, _, err := connectToWebSocket(boykisserToken)
	require.NoError(t, err, "boykisser failed to connect to the web socket")
	defer boykisserWS.Close()

	var messageID string

	t.Run("cutie sends message", func(t *testing.T) {
		message := "Haiiii! >w<"
		resp, _ := sendTestMessage(GatewayURL, cutieToken, message)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		messageID = respBody["message_id"]

		//wait for ws
		event := waitForEvent(t, boykisserWS, "userMessage", 5*time.Second)

		//extract payload
		payload := event["payload"].(map[string]any)
		assert.Equal(t, message, payload["content"])
		assert.Equal(t, user1, payload["username"])
	})

	t.Run("cutie edits message", func(t *testing.T) {
		newContent := "sowwy (edited)"
		resp := sendEditRequest(GatewayURL, cutieToken, messageID, newContent)
		require.Equal(t, http.StatusOK, resp.StatusCode, "failed to edit the message :c (tip: check docker logs)")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		event := waitForEvent(t, boykisserWS, "editMessage", 5*time.Second)
		payload := event["payload"].(map[string]any)
		assert.Equal(t, newContent, payload["content"])
		assert.Equal(t, messageID, payload["message_id"])
	})

	t.Run("boykisser tries to delete cutie's message (Security Check)", func(t *testing.T) {
		resp := sendDeleteRequest(GatewayURL, boykisserToken, messageID)

		assert.Equal(t, http.StatusForbidden, resp.StatusCode, "")
	})

	t.Run("cutie deletes message", func(t *testing.T) {
		resp := sendDeleteRequest(GatewayURL, cutieToken, messageID)
		require.Equal(t, http.StatusOK, resp.StatusCode, "failed to edit the message :c (tip: check docker logs)")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		event := waitForEvent(t, boykisserWS, "deleteMessage", 5*time.Second)
		payload := event["payload"].(map[string]any)
		assert.Equal(t, messageID, payload["message_id"])
	})
	t.Run("boykisser sends message and cutie reads history", func(t *testing.T) {
		message := "Haiii! >w<"
		resp, _ := sendTestMessage(GatewayURL, boykisserToken, message)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		messageID = respBody["message_id"]

		resp, _ = sendGetHistoryRequest(GatewayURL, cutieToken, 10)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var messagesResp map[string]any
		json.NewDecoder(resp.Body).Decode(&messagesResp)
		messages := messagesResp["messages"].([]any)
		assert.Equal(t, 1, len(messages))

		returnedMessage := messages[0].(map[string]any)

		assert.Equal(t, messageID, returnedMessage["message_id"])
		assert.Equal(t, message, returnedMessage["content"])
	})
}

func sendGetHistoryRequest(url, token string, limit int) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url+"/messages?limit="+strconv.Itoa(limit), nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	req.Header.Add("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func sendTestMessage(url, token, content string) (*http.Response, error) {
	body, _ := json.Marshal(map[string]string{
		"content": content,
	})

	req, err := http.NewRequest(http.MethodPost, url+"/messages", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	req.Header.Add("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func connectToWebSocket(token string) (*websocket.Conn, *http.Response, error) {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+token)
	dialer := websocket.DefaultDialer
	return dialer.Dial(WSURL, header)
}

func loginAndGetToken(username, password string) (string, error) {
	//password field is for future
	body, _ := json.Marshal(loginRequest{
		Username: username,
	})

	resp, err := http.Post(GatewayURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to send POST request to /login")
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("login response status code should be 200")
	}

	defer resp.Body.Close()

	var loginResp loginResponse

	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("login response body doesn't match")
	}

	return loginResp.Token, nil
}

func waitForEvent(t *testing.T, ws *websocket.Conn, expectedType string, timeout time.Duration) map[string]any {
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			t.Fatalf("Timeout: Nie otrzymano eventu '%s' w ciągu %v", expectedType, timeout)
			return nil
		}

		ws.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

		var msg map[string]any
		err := ws.ReadJSON(&msg)

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}

			t.Fatalf("Fatal error QwQ: %v. Disconnecting from web socket", err)
			return nil
		}

		t.Logf("received ws event: %v", msg)

		if eventType, ok := msg["type"].(string); ok && eventType == expectedType {
			return msg
		}
	}
}

func sendEditRequest(url, token, msgID, content string) *http.Response {
	body, _ := json.Marshal(map[string]string{
		"message_id": msgID,
		"content":    content,
	})
	req, _ := http.NewRequest(http.MethodPatch, url+"/messages", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, _ := http.DefaultClient.Do(req)
	return resp
}

func sendDeleteRequest(url, token, msgID string) *http.Response {
	body, _ := json.Marshal(map[string]string{
		"message_id": msgID,
	})
	req, _ := http.NewRequest(http.MethodDelete, url+"/messages", bytes.NewBuffer(body))
	req.Header.Add("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	return resp
}
