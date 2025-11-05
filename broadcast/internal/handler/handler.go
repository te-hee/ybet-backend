package handler

import (
	"broadcast/config"
	"broadcast/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	mutex          *sync.Mutex
	conns          map[string]*websocket.Conn
	messageChannel chan models.Message
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebsocketHandler(msgChan chan models.Message) *WebSocketHandler {
	return &WebSocketHandler{
		mutex:          &sync.Mutex{},
		conns:          make(map[string]*websocket.Conn),
		messageChannel: msgChan,
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Error encoding response:", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, models.ErrorResponse{
		Error: message,
	})
}

func (websockethandler *WebSocketHandler) WsHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserClaims
	var err error
	if !*config.NoAuth {
		token := r.URL.Query().Get("token")

		if token == "" {
			writeError(w, http.StatusUnauthorized, "provide token in url param")
		}
		user, err = verifyJWT()
		if err != nil {
			writeError(w, 401, "error verifying JWT token :c")
		}

	} else {
		user.Uuid = uuid.New()
		user.Username = fmt.Sprintf("user%d", len(websockethandler.conns))
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//Tymon :3
		writeError(w, 418, "I'm a teapot")
	}
	websockethandler.mutex.Lock()

	websockethandler.conns[user.Uuid.String()] = conn

	websockethandler.mutex.Unlock()

	websockethandler.messageChannel <- models.Message{
		Type: models.UserListUpdateType,
		Payload: models.UserListUpdate{
			Action: models.Connect,
			Uuid:   user.Uuid.String(),
		},
	}

	websockethandler.messageChannel <- models.Message{
		Type: models.SystemMessageType,
		Payload: models.SystemMessage{
			Content: fmt.Sprintf("%s joined! :3", user.Username),
		},
	}

}

func (websockethandler *WebSocketHandler) BroadcastMessages() {
	for {
		msg := <-websockethandler.messageChannel
		log.Println("got message from channel")
		websockethandler.mutex.Lock()

		for user, conn := range websockethandler.conns {
			log.Println("sending message to user")
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("cant send message to user")
				conn.Close()
				delete(websockethandler.conns, user)
			}
		}

		websockethandler.mutex.Unlock()
		log.Println("sent message to all cons")

	}
}

func verifyJWT() (models.UserClaims, error) {
	return models.UserClaims{}, nil
}
