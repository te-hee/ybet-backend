package models

import messagev1 "backend/proto/message/v1"

type Message struct {
	Id        string
	Message   string
	Timestamp int64
}

func (m Message) ToProto() *messagev1.Message {
	return &messagev1.Message{
		Id:        m.Id,
		Content:   m.Message,
		Timestamp: m.Timestamp,
	}
}
