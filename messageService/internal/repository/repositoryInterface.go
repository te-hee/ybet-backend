package repository

import "backend/messageService/internal/models"

type Repository interface {
	SaveMessage(message models.Message)
	GetMessages(limit int) []models.Message
}
