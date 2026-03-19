package service

import (
	"roomService/internal/ports"
)

type roomService struct {
	repo           ports.RoomRepository
	eventPublisher ports.EventPublisher
	keyRepo        ports.KeyRepository
}

type Config struct {
	Repo      ports.RoomRepository
	Publisher ports.EventPublisher
	KeyRepo   ports.KeyRepository
}

func New(cfg Config) ports.RoomService {
	return &roomService{
		repo:           cfg.Repo,
		eventPublisher: cfg.Publisher,
		keyRepo:        cfg.KeyRepo,
	}
}
