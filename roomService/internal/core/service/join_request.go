package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *roomService) CreateJoinRequest(ctx context.Context, roomUUID string, publicKey string) error {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	username, err := contextkeys.UsernameFromContext(ctx)
	if err != nil {
		return err
	}

	joinRequest := domain.JoinRequest{
		RoomUUID:  roomUUID,
		UserUUID:  userID,
		Username:  username,
		PublicKey: publicKey,
		Status:    domain.RequestStatusPending,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateJoinRequest(ctx, joinRequest)
}

func (s *roomService) GetJoinRequests(ctx context.Context, roomUUID string) ([]domain.JoinRequest, error) {
	var joinRequests []domain.JoinRequest

	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return joinRequests, err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, userID, roomUUID)
	if err != nil {
		return joinRequests, err
	}
	if !isAdmin {
		return joinRequests, status.Errorf(codes.PermissionDenied, "Only admins can get join requests")
	}

	return s.repo.GetJoinRequests(ctx, roomUUID)
}

func (s *roomService) RespondToJoinRequest(ctx context.Context, roomUUID string, userUUID string, decision domain.RequestStatus) error {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, userID, roomUUID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return status.Errorf(codes.PermissionDenied, "Only admins can respond to join requests")
	}

	return s.repo.UpdateJoinRequestStatus(ctx, roomUUID, userUUID, decision)
}
