package service

import (
	"backend/messageService/internal/models"
	"backend/messageService/internal/repository"
)

type ServiceLayer struct {
	repo repository.Repository
}

func New(repo repository.Repository) *ServiceLayer {
	return &ServiceLayer{
		repo: repo,
	}
}

func (s *ServiceLayer) SaveMessage(message models.Message) {
	s.repo.SaveMessage(message)
}

func (s *ServiceLayer) GetMessages(limit int) []models.Message {
	return s.repo.GetMessages(limit)
}
