package repository

import (
	"backend/messageService/internal/models"
)

type MemoryRepo struct {
	messages []models.Message
}

func NewInMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		messages: make([]models.Message, 0),
	}
}

func (r *MemoryRepo) SaveMessage(message models.Message) {
	r.messages = append(r.messages, message)
}
func (r *MemoryRepo) GetMessages(limit int) []models.Message {
	length := len(r.messages)
	if limit > length {
		limit = length
	}
	return r.messages[length-limit : length]
}
