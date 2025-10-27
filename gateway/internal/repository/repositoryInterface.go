package repository

import(
	"backend/gateway/internal/model"
)

type Repository interface{
	GetMessageHistory(limit uint32) ([]model.Message, error)
	SendMessage(message model.Message) (error)
}