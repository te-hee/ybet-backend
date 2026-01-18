package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sendMessageConfig struct {
	expectedCode int
	withToken    bool
}

type loginConfig struct {
	expectedCode int
}

func TestLogin(t *testing.T) {
	loginResp, err := loginUser("cutie1")
	if err != nil {
		t.Errorf("login error :c\n %v", err)
		return
	}
	assert.NotEmpty(t, loginResp.Token)
	assert.NotEmpty(t, loginResp.UserId)
}

func TestLoginNoUsername(t *testing.T) {
	_, err := loginUser("", loginConfig{
		expectedCode: http.StatusBadRequest,
	})
	if err != nil {
		t.Errorf("login error :c\n %v", err)
		return
	}
}
func TestSendMessage_Table(t *testing.T) {
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

func loginUser(username string, cfg ...loginConfig) (*loginResponse, error) {
	expectedCode := http.StatusOK
	if len(cfg) > 0 {
		expectedCode = cfg[0].expectedCode
	}
	body, _ := json.Marshal(loginRequest{
		Username: username,
	})

	resp, err := http.Post(GatewayURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request to /login")
	}
	if resp.StatusCode != expectedCode {
		return nil, fmt.Errorf("login response status code should be 200")
	}

	defer resp.Body.Close()

	var loginResp loginResponse

	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, fmt.Errorf("login response body doesn't match")
	}

	return &loginResp, nil
}
