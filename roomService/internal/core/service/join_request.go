package service

import (
	"context"
	"roomService/internal/core/domain"
)

func (s *roomService) CreateJoinRequest(ctx context.Context, roomUUID string, publicKey string) error {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) GetJoinRequests(ctx context.Context, roomUUID string) ([]domain.JoinRequest, error) {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) RespondToJoinRequest(ctx context.Context, roomUUID string, userUUID string, decision domain.RequestStatus) error {
	// TODO: implement
	panic("not implemented")
}
