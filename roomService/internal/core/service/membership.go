package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *roomService) GetRoomMembers(ctx context.Context, roomUUID string) ([]domain.RoomMember, error) {
	var members []domain.RoomMember
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return members, err
	}

	isMember, err := s.repo.CheckIsMember(ctx, userID, roomUUID)
	if err != nil {
		return members, err
	}

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
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	room, err := s.repo.GetRoom(ctx, roomUUID)
	if err != nil {
		return err
	}
	if room.AdminID == userID {
		return status.Error(codes.PermissionDenied, "Admins cannot leave the room. Please delete the room instead")
	}

	err = s.repo.RemoveMember(ctx, roomUUID, userID)
	if err != nil {
		return err
	}

	return s.eventPublisher.PublishMemberLeft(ctx, roomUUID, userID)
}

func (s *roomService) RemoveMember(ctx context.Context, roomUUID string, userUUID string) error {
	requesterID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, requesterID, roomUUID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return status.Errorf(codes.PermissionDenied, "Only admins can remove members")
	}

	room, err := s.repo.GetRoom(ctx, roomUUID)
	if err != nil {
		return err
	}
	if room.AdminID == userUUID {
		return status.Error(codes.PermissionDenied, "Admins cannot be removed from the room")
	}

	err = s.repo.RemoveMember(ctx, roomUUID, userUUID)
	if err != nil {
		return err
	}

	return s.eventPublisher.PublishMemberLeft(ctx, roomUUID, userUUID)
}
