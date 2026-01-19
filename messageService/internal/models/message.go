package models

import (
	messagev1 "backend/proto/message/v1"

	"github.com/google/uuid"
)

type Message struct {
	Username  string    `json:"username,omitempty"`
	UserId    uuid.UUID `json:"user_id,omitempty"`
	Id        uuid.UUID `json:"message_id,omitempty"`
	Message   string    `json:"content,omitempty"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

func (m Message) ToProto() *messagev1.Message {
	return &messagev1.Message{
		Username:  m.Username,
		UserId:    m.UserId.String(),
		Uuid:      m.Id.String(),
		Content:   m.Message,
		Timestamp: m.Timestamp,
	}
}
