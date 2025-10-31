package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	UserId uuid.UUID
	Conn   *websocket.Conn
}
