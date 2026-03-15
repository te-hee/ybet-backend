package service

import (
	"context"
	"roomService/internal/core/domain"
)

func (s *roomService) CreateRoom(ctx context.Context, name string, isPrivate bool, groupID string) (domain.CreateRoomResult, error) {
	// TODO: implement
	panic("not implemented")
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
