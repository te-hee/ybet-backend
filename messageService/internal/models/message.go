package models

import (
	messagev2 "backend/proto/message/v2"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	Username  string    `json:"username,omitempty"`
	UserId    uuid.UUID `json:"user_id,omitempty"`
	Id        uuid.UUID `json:"message_id,omitempty"`
	Message   string    `json:"content,omitempty"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

func (m Message) ToProto() *messagev2.Message {
	return &messagev2.Message{
		Username:  m.Username,
		UserId:    m.UserId.String(),
		MessageId: m.Id.String(),
		Content:   m.Message,
		CreatedAt: timestamppb.New(time.Unix(m.Timestamp, 0)),
	}
}
