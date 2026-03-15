package service

import "roomService/internal/ports"

type roomService struct {
	repo      ports.RoomRepository
	publisher ports.EventPublisher
}

func New(repo ports.RoomRepository, publisher ports.EventPublisher) ports.RoomService {
	return &roomService{
		repo:      repo,
		publisher: publisher,
	}
}
