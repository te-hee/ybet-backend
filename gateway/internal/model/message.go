package model

import (
	v1 "backend/proto/message/v1"

	"github.com/google/uuid"
)

type Message struct {
	Uuid      uuid.UUID `json:"uuid"`
	Content   string    `json:"content"`
	Timestamp uint64    `json:"timestamp"`
	UserId    uuid.UUID `json:"userId"`
}

func ProtoToModel(protoMessage *v1.Message) (*Message, error) {
	id, err := uuid.Parse(protoMessage.Uuid)

	if err != nil {
		return nil, err
	}
	userId, err := uuid.Parse(protoMessage.Uuid)
	if err != nil {
		return nil, err
	}

	message := &Message{
		Uuid:      id,
		Content:   protoMessage.Content,
		Timestamp: uint64(protoMessage.Timestamp),
		UserId:    userId,
	}
	return message, nil
}

func ModelToProto(message Message) *v1.Message {
	protoMessage := &v1.Message{
		Uuid:      message.Uuid.String(),
		Content:   message.Content,
		Timestamp: int64(message.Timestamp),
		UserId:    message.UserId.String(),
	}
	return protoMessage
}
