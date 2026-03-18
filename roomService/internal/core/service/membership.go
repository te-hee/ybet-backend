package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
	"slices"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *roomService) GetRoomMembers(ctx context.Context, roomUUID string) ([]domain.RoomMember, error) {
	var members []domain.RoomMember
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return members, err
	}

	userRooms, err := s.repo.GetUserRooms(ctx, userID)
	if err != nil {
		return members, err
	}

	isMember := slices.ContainsFunc(userRooms, func(room domain.UserRoom) bool {
		return room.RoomUUID == roomUUID
	})

	if !isMember {
		return members, status.Error(codes.PermissionDenied, "Only room members can check other members")
	}

	members, err = s.repo.GetRoomMembers(ctx, roomUUID)
	if err != nil {
		return members, err
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
