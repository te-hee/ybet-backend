package handler

import (
	"backend/broadcast/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	mutex          *sync.Mutex
	conns          map[models.User]bool
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
		conns:          make(map[models.User]bool),
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//Tymon :3
		writeError(w, 418, "I'm a teapot")
	}
	user := models.User{
		UserId: uuid.New(),
		Conn:   conn,
	}
	websockethandler.mutex.Lock()

	websockethandler.conns[user] = true

	websockethandler.mutex.Unlock()

}

func (websockethandler *WebSocketHandler) BroadcastMessages() {
	for {
		msg := <-websockethandler.messageChannel
		log.Println("got message from channel")
		websockethandler.mutex.Lock()

		for user := range websockethandler.conns {
			log.Println("sending message to user")
			if err := user.Conn.WriteJSON(msg); err != nil {
				log.Println("cant send message to user")
				user.Conn.Close()
				delete(websockethandler.conns, user)
			}
		}

		websockethandler.mutex.Unlock()
		log.Println("sent message to all cons")

	}
}
