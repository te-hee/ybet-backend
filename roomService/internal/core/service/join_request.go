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

	_, err = s.repo.GetRoom(ctx, roomUUID)
	if err != nil {
		return err
	}

	isMember, err := s.repo.CheckIsMember(ctx, userID, roomUUID)
	if err != nil {
		return err
	}
	if isMember {
		return status.Error(codes.AlreadyExists, "User is already a member of this room")
	}

	existingReq, err := s.repo.GetJoinRequest(ctx, roomUUID, userID)
	if err == nil && existingReq.Status == domain.RequestStatusPending {
		return status.Error(codes.AlreadyExists, "A pending join request already exists")
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

	err = s.repo.UpdateJoinRequestStatus(ctx, roomUUID, userUUID, decision)
	if err != nil {
		return err
	}

	if decision == domain.RequestStatusAccepted {
		err = s.repo.AddMember(ctx, roomUUID, userUUID)
		if err != nil {
			s.repo.UpdateJoinRequestStatus(ctx, roomUUID, userUUID, domain.RequestStatusPending)
			return err
		}

		err = s.keyRepo.SaveKey(ctx, domain.PendingKey{
			RoomUUID: roomUUID,
			UserUUID: userUUID,
		})
		if err != nil {
			s.repo.RemoveMember(ctx, roomUUID, userUUID)
			s.repo.UpdateJoinRequestStatus(ctx, roomUUID, userUUID, domain.RequestStatusPending)
			return err
		}

		return s.eventPublisher.PublishMemberJoined(ctx, roomUUID, userUUID)
	}
	return nil
}

func (s *roomService) GetPendingKeys(ctx context.Context) ([]domain.PendingKey, error) {
	userId, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.keyRepo.GetUserKeys(ctx, userId)
}

func (s *roomService) AcknowledgeKeyDelivery(ctx context.Context, roomUUID string) error {
	userId, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	return s.keyRepo.DeleteKey(ctx, roomUUID, userId)
}
