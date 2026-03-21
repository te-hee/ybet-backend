package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

func (s *roomService) CreateInvite(ctx context.Context, req domain.CreateInviteDTO) (domain.CreateInviteResult, error) {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return domain.CreateInviteResult{}, err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, userID, req.RoomUUID)
	if err != nil {
		return domain.CreateInviteResult{}, err
	}
	if !isAdmin {
		return domain.CreateInviteResult{}, domain.NewError(domain.CodePermissionDenied, "user is not an admin")
	}

	if req.UsesLeft < 0 {
		return domain.CreateInviteResult{}, domain.NewError(domain.CodeInvalidArgument, "usesLeft cannot be negative")
	}

	var expirationTime time.Time
	if req.ExpiresAt != nil {
		expiresAt := req.ExpiresAt.AsTime().Unix()
		if expiresAt > 0 {
			expirationTime = time.Unix(expiresAt, 0)
			if expirationTime.Before(time.Now()) {
				return domain.CreateInviteResult{}, domain.NewError(domain.CodeInvalidArgument, "expiresAt cannot be in the past")
			}
		}
	}

	inviteID := uuid.New().String()

	invite := domain.RoomInvite{
		InviteID:  inviteID,
		RoomUUID:  req.RoomUUID,
		UsesLeft:  req.UsesLeft,
		ExpiresAt: expirationTime,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateInvite(ctx, invite); err != nil {
		return domain.CreateInviteResult{}, err
	}

	return domain.CreateInviteResult{InviteID: inviteID}, nil
}

func (s *roomService) GetInvite(ctx context.Context, req domain.GetInviteDTO) (domain.RoomInvite, error) {
	return s.repo.GetInvite(ctx, req.InviteID)
}

func (s *roomService) DeleteInvite(ctx context.Context, req domain.DeleteInviteDTO) error {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	isAdmin, err := s.repo.CheckIsAdmin(ctx, userID, req.RoomUUID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return domain.NewError(domain.CodePermissionDenied, "user is not an admin")
	}

	return s.repo.DeleteInvite(ctx, req.InviteID, req.RoomUUID)
}

func (s *roomService) JoinViaInvite(ctx context.Context, req domain.JoinViaInviteDTO) error {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	invite, err := s.repo.GetInvite(ctx, req.InviteID)
	if err != nil {
		return err
	}

	if invite.UsesLeft != 0 && invite.UsesLeft <= 0 {
		return domain.NewError(domain.CodeInvalidArgument, "invite is fully consumed")
	}

	if !invite.ExpiresAt.IsZero() && invite.ExpiresAt.Before(time.Now()) {
		return domain.NewError(domain.CodeInvalidArgument, "invite is expired")
	}

	isMember, err := s.repo.CheckIsMember(ctx, userID, invite.RoomUUID)
	if err != nil {
		return err
	}
	if isMember {
		return domain.NewError(domain.CodeAlreadyExists, "user is already a member")
	}

	//TODO: key exchange mechanism

	if invite.UsesLeft > 0 {
		err = s.repo.DecrementInviteUses(ctx, req.InviteID)
		if err != nil {
			return err
		}
	}

	return s.eventPublisher.PublishMemberJoined(ctx, invite.RoomUUID, userID)
}
