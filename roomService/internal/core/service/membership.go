package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
)

func (s *roomService) GetRoomMembers(ctx context.Context, roomUUID string) ([]domain.RoomMember, error) {
	var members []domain.RoomMember
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return members, err
	}

	members, err = s.repo.GetRoomMembers(ctx, roomUUID)
	if err != nil {

	}
	return members, nil
}

func (s *roomService) LeaveRoom(ctx context.Context, roomUUID string) error {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) RemoveMember(ctx context.Context, roomUUID string, userUUID string) error {
	// TODO: implement
	panic("not implemented")
}
