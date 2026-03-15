package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func (s *roomService) CreateRoom(ctx context.Context, name string, isPrivate bool, groupID string) (domain.CreateRoomResult, error) {
	var result domain.CreateRoomResult

	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return result, err
	}

	roomID := uuid.New().String()
	createdAt := time.Now()
	room := domain.Room{
		RoomUUID:    roomID,
		Name:        name,
		AdminID:     userID,
		IsPrivate:   isPrivate,
		GroupID:     groupID,
		MemberCount: 1,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	if err = s.repo.CreateRoom(ctx, room); err != nil {
		return result, err
	}

	result = domain.CreateRoomResult{
		RoomUUID:  roomID,
		CreatedAt: createdAt,
	}

	return result, nil
}

func (s *roomService) GetRoom(ctx context.Context, roomUUID string) (domain.Room, error) {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) UpdateRoomName(ctx context.Context, roomUUID string, name string) error {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) DeleteRoom(ctx context.Context, roomUUID string) error {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) GetUserRooms(ctx context.Context) ([]domain.UserRoom, error) {
	// TODO: implement
	panic("not implemented")
}
