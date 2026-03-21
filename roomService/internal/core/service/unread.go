package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
)

func (s *roomService) MarkAsRead(ctx context.Context, req domain.MarkAsReadDTO) error {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repo.MarkAsRead(ctx, req.RoomUUID, userID, req.LastReadMessageID)
}

func (s *roomService) GetUnreadCount(ctx context.Context, req domain.GetUnreadCountDTO) (domain.GetUnreadCountResult, error) {
	userID, err := contextkeys.UserUUIDFromContext(ctx)
	if err != nil {
		return domain.GetUnreadCountResult{}, err
	}

	count, err := s.repo.GetUnreadCount(ctx, req.RoomUUID, userID)
	if err != nil {
		return domain.GetUnreadCountResult{}, err
	}

	return domain.GetUnreadCountResult{UnreadCount: count}, nil
}
