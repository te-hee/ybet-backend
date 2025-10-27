package models

import (
	messagev1 "backend/proto/message/v1"

	"github.com/google/uuid"
)

type Message struct {
	Id        uuid.UUID
	Message   string
	Timestamp int64
}

func (m Message) ToProto() *messagev1.Message {
	return &messagev1.Message{
		Uuid:      m.Id.String(),
		Content:   m.Message,
		Timestamp: m.Timestamp,
	}
}
