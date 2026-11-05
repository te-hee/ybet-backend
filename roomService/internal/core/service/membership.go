package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
)

func (s *roomService) GetRoomMembers(ctx context.Context, req domain.GetRoomMembersDTO) ([]domain.RoomMember, error) {
	var members []domain.RoomMember
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return members, err
	}

	isMember, err := s.repo.CheckIsMember(ctx, userID, req.RoomUUID)
	if err != nil {
		return members, err
	}

	if !isMember {
		return members, domain.NewError(domain.CodePermissionDenied, "Only room members can check other members")
	}

	members, err = s.repo.GetRoomMembers(ctx, req.RoomUUID)
	if err != nil {
		return members, err
	}
	return members, nil
}

func (s *roomService) LeaveRoom(ctx context.Context, req domain.LeaveRoomDTO) error {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	room, err := s.repo.GetRoom(ctx, req.RoomUUID)
	if err != nil {
		return err
	}
	if room.AdminID == userID {
		return domain.NewError(domain.CodePermissionDenied, "Admins cannot leave the room. Please delete the room instead")
	}

	err = s.repo.RemoveMember(ctx, req.RoomUUID, userID)
	if err != nil {
		return err
	}

	return s.eventPublisher.PublishMemberLeft(ctx, req.RoomUUID, userID)
}

func (s *roomService) RemoveMember(ctx context.Context, req domain.RemoveMemberDTO) error {
	requesterID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, requesterID, req.RoomUUID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return domain.NewError(domain.CodePermissionDenied, "Only admins can remove members")
	}

	room, err := s.repo.GetRoom(ctx, req.RoomUUID)
	if err != nil {
		return err
	}
	if room.AdminID == req.UserUUID {
		return domain.NewError(domain.CodePermissionDenied, "Admins cannot be removed from the room")
	}

	err = s.repo.RemoveMember(ctx, req.RoomUUID, req.UserUUID)
	if err != nil {
		return err
	}

	return s.eventPublisher.PublishMemberLeft(ctx, req.RoomUUID, req.UserUUID)
}
