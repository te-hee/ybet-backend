package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *roomService) CreateRoom(ctx context.Context, name string, isPrivate bool, groupID string) (domain.CreateRoomResult, error) {
	var result domain.CreateRoomResult

	if name == "" {
		return result, status.Error(codes.InvalidArgument, "room name cannot be empty")
	}

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

func (s *roomService) GetRoom(ctx context.Context, roomUUID string) (domain.Room, error) {
	room, err := s.repo.GetRoom(ctx, roomUUID)
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

	isMember, err := s.repo.CheckIsMember(ctx, userID, roomUUID)
	if err != nil {
		return domain.Room{}, err
	}

	if !isMember {
		return domain.Room{}, status.Errorf(codes.PermissionDenied, "Only room members can access private rooms")
	}

	return room, nil
}

func (s *roomService) UpdateRoomName(ctx context.Context, roomUUID string, name string) error {
	if name == "" {
		return status.Error(codes.InvalidArgument, "room name cannot be empty")
	}

	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, userID, roomUUID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return status.Errorf(codes.PermissionDenied, "Only admins can update room name")
	}

	return s.repo.UpdateRoomName(ctx, roomUUID, name)
}

func (s *roomService) DeleteRoom(ctx context.Context, roomUUID string) error {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, userID, roomUUID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return status.Errorf(codes.PermissionDenied, "Only admins can delete rooms")
	}

	err = s.repo.DeleteRoom(ctx, roomUUID)
	if err != nil {
		return err
	}

	return s.eventPublisher.PublishRoomDeleted(ctx, roomUUID)
}

func (s *roomService) GetUserRooms(ctx context.Context) ([]domain.UserRoom, error) {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserRooms(ctx, userID)
}
