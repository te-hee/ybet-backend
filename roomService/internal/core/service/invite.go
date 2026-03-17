package service

import (
	"context"
	"roomService/internal/core/domain"
)

func (s *roomService) CreateInvite(ctx context.Context, roomUUID string, usesLeft int32, expiresAt int64) (domain.CreateInviteResult, error) {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) GetInvite(ctx context.Context, inviteID string) (domain.RoomInvite, error) {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) DeleteInvite(ctx context.Context, inviteID string, roomUUID string) error {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) JoinViaInvite(ctx context.Context, inviteID string) error {
	// TODO: implement
	panic("not implemented")
}
