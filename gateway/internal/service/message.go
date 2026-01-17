package service

import (
	"gateway/internal/model"
	"gateway/internal/repository"
)

type MessageService struct {
	repo *repository.RepositoryGrpc
}

func NewMessageService(repo *repository.RepositoryGrpc) *MessageService {
	messageService := &MessageService{
		repo: repo,
	}
	return messageService
}

func (s *MessageService) GetMessageHistory(limit uint32) ([]model.Message, error) {
	users, err := s.repo.GetMessageHistory(limit)
	return users, err
}

func (s *MessageService) SendMessage(message model.InputMessage) error {
	err := s.repo.SendMessage(message)
	return err
}

func (s *MessageService) EditMessage(editRequest model.EditMessageRequest) error {
	err := s.repo.EditMessage(editRequest)
	return err
}

func (s *MessageService) DeleteMessage(deleteRequest model.DeleteMessageRequest) error {
	err := s.repo.DeleteMessage(deleteRequest)
	return err
}
