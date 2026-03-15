package service

import (
	"context"
	"roomService/internal/core/domain"
)

func (s *roomService) GetRoomMembers(ctx context.Context, roomUUID string) ([]domain.RoomMember, error) {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) LeaveRoom(ctx context.Context, roomUUID string) error {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) RemoveMember(ctx context.Context, roomUUID string, userUUID string) error {
	// TODO: implement
	panic("not implemented")
}
