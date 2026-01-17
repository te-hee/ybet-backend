package repository

import "messageService/internal/models"

type Repository interface {
	SaveMessage(message models.Message)
	GetMessages(limit int) []models.Message
	EditMessage(editMessage models.EditMessage) error
	DeleteMessage(deleteMssage models.DeleteMessage) error
}
