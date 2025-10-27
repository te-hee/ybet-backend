package service

import (
	"backend/gateway/internal/repository"
	"backend/gateway/internal/model"
)

type MessageService struct{
	repo *repository.RepositoryGrpc
}

func NewMessageService(repo *repository.RepositoryGrpc) *MessageService{
	messageService := &MessageService{
		repo: repo,
	}
	return  messageService
}

func (s* MessageService) GetMessageHistory(limit uint32) ([]model.Message, error){
	users, err := s.repo.GetMessageHistory(limit);
	return users, err
}

func (s* MessageService) SendMessage(content string) (error){
	err := s.repo.SendMessage(content)
	return  err
}