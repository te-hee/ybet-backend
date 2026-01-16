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
