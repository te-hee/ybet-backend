package service

import (
	"gateway/internal/client"
	"gateway/internal/model"
)

type MessageService struct {
	client client.MessageClient
}

func NewMessageService(client client.MessageClient) *MessageService {
	messageService := &MessageService{
		client: client,
	}
	return messageService
}

func (s *MessageService) GetMessageHistory(limit uint32) ([]model.Message, error) {
	users, err := s.client.GetMessageHistory(limit)
	return users, err
}

func (s *MessageService) SendMessage(message model.InputMessage) (*model.OutputSendMessege, error) {
	resp, err := s.client.SendMessage(message)
	return resp, err
}

func (s *MessageService) EditMessage(editRequest model.EditMessageRequest) error {
	err := s.client.EditMessage(editRequest)
	return err
}

func (s *MessageService) DeleteMessage(deleteRequest model.DeleteMessageRequest) error {
	err := s.client.DeleteMessage(deleteRequest)
	return err
}
