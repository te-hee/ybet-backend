package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
	"time"
)

func (s *roomService) CreateJoinRequest(ctx context.Context, req domain.CreateJoinRequestDTO) error {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	username, err := contextkeys.UsernameFromContext(ctx)
	if err != nil {
		return err
	}

	_, err = s.repo.GetRoom(ctx, req.RoomUUID)
	if err != nil {
		return err
	}

	isMember, err := s.repo.CheckIsMember(ctx, userID, req.RoomUUID)
	if err != nil {
		return err
	}
	if isMember {
		return domain.NewError(domain.CodeAlreadyExists, "User is already a member of this room")
	}

	existingReq, err := s.repo.GetJoinRequest(ctx, req.RoomUUID, userID)
	if err == nil && existingReq.Status == domain.RequestStatusPending {
		return domain.NewError(domain.CodeAlreadyExists, "A pending join request already exists")
	}

	joinRequest := domain.JoinRequest{
		RoomUUID:  req.RoomUUID,
		UserUUID:  userID,
		Username:  username,
		PublicKey: req.PublicKey,
		Status:    domain.RequestStatusPending,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateJoinRequest(ctx, joinRequest)
}

func (s *roomService) GetJoinRequests(ctx context.Context, req domain.GetJoinRequestsDTO) ([]domain.JoinRequest, error) {
	var joinRequests []domain.JoinRequest

	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return joinRequests, err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, userID, req.RoomUUID)
	if err != nil {
		return joinRequests, err
	}
	if !isAdmin {
		return joinRequests, domain.NewError(domain.CodePermissionDenied, "Only admins can get join requests")
	}

	return s.repo.GetJoinRequests(ctx, req.RoomUUID)
}

func (s *roomService) RespondToJoinRequest(ctx context.Context, req domain.RespondToJoinRequestDTO) error {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, userID, req.RoomUUID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return domain.NewError(domain.CodePermissionDenied, "Only admins can respond to join requests")
	}

	err = s.repo.UpdateJoinRequestStatus(ctx, req.RoomUUID, req.UserUUID, req.Decision)
	if err != nil {
		return err
	}

	if req.Decision == domain.RequestStatusAccepted {
		err = s.repo.AddMember(ctx, req.RoomUUID, req.UserUUID)
		if err != nil {
			s.repo.UpdateJoinRequestStatus(ctx, req.RoomUUID, req.UserUUID, domain.RequestStatusPending)
			return err
		}

		err = s.keyRepo.SaveKey(ctx, domain.PendingKey{
			RoomUUID: req.RoomUUID,
			UserUUID: req.UserUUID,
		})
		if err != nil {
			s.repo.RemoveMember(ctx, req.RoomUUID, req.UserUUID)
			s.repo.UpdateJoinRequestStatus(ctx, req.RoomUUID, req.UserUUID, domain.RequestStatusPending)
			return err
		}

		return s.eventPublisher.PublishMemberJoined(ctx, req.RoomUUID, req.UserUUID)
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

func (s *roomService) AcknowledgeKeyDelivery(ctx context.Context, req domain.AcknowledgeKeyDeliveryDTO) error {
	userId, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	return s.keyRepo.DeleteKey(ctx, req.RoomUUID, userId)
}
