package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func (s *roomService) CreateRoom(ctx context.Context, req domain.CreateRoomDTO) (domain.CreateRoomResult, error) {
	var result domain.CreateRoomResult

	if req.Name == "" {
		return result, domain.NewError(domain.CodeInvalidArgument, "room name cannot be empty")
	}

	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return result, err
	}

	roomID := uuid.New().String()
	createdAt := time.Now()
	room := domain.Room{
		RoomUUID:    roomID,
		Name:        req.Name,
		AdminID:     userID,
		IsPrivate:   req.IsPrivate,
		GroupID:     req.GroupID,
		MemberCount: 1,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	if err = s.repo.CreateRoom(ctx, room); err != nil {
		return result, err
	}

	err = s.eventPublisher.PublishRoomCreated(ctx, room)
	if err != nil {
		return result, err
	}

	result = domain.CreateRoomResult{
		RoomUUID:  roomID,
		CreatedAt: createdAt,
	}

	return result, nil
}

func (s *roomService) GetRoom(ctx context.Context, req domain.GetRoomDTO) (domain.Room, error) {
	room, err := s.repo.GetRoom(ctx, req.RoomUUID)
	if err != nil {
		return domain.Room{}, err
	}

	if !room.IsPrivate {
		return room, nil
	}

	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return domain.Room{}, err
	}

	isMember, err := s.repo.CheckIsMember(ctx, userID, req.RoomUUID)
	if err != nil {
		return domain.Room{}, err
	}

	if !isMember {
		return domain.Room{}, domain.NewError(domain.CodePermissionDenied, "Only room members can access private rooms")
	}

	return room, nil
}

func (s *roomService) UpdateRoomName(ctx context.Context, req domain.UpdateRoomNameDTO) error {
	if req.Name == "" {
		return domain.NewError(domain.CodeInvalidArgument, "room name cannot be empty")
	}

	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, userID, req.RoomUUID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return domain.NewError(domain.CodePermissionDenied, "Only admins can update room name")
	}

	return s.repo.UpdateRoomName(ctx, req.RoomUUID, req.Name)
}

func (s *roomService) DeleteRoom(ctx context.Context, req domain.DeleteRoomDTO) error {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, userID, req.RoomUUID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return domain.NewError(domain.CodePermissionDenied, "Only admins can delete rooms")
	}

	err = s.repo.DeleteRoom(ctx, req.RoomUUID)
	if err != nil {
		return err
	}

	return s.eventPublisher.PublishRoomDeleted(ctx, req.RoomUUID)
}

func (s *roomService) GetUserRooms(ctx context.Context) ([]domain.UserRoom, error) {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserRooms(ctx, userID)
}
