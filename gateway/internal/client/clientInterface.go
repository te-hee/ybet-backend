package client

import "gateway/internal/model"

type MessageClient interface {
	GetMessageHistory(limit uint32) ([]model.Message, error)
	SendMessage(message model.InputMessage) (*model.OutputSendMessege, error)
	EditMessage(editRequest model.EditMessageRequest) error
	DeleteMessage(deleteRequest model.DeleteMessageRequest) error
}
