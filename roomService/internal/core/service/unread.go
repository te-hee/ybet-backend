package service

import (
	"context"
	"roomService/internal/core/domain"
)

func (s *roomService) MarkAsRead(ctx context.Context, roomUUID string, lastReadMessageID string) error {
	// TODO: implement
	panic("not implemented")
}

func (s *roomService) GetUnreadCount(ctx context.Context, roomUUID string) (domain.GetUnreadCountResult, error) {
	// TODO: implement
	panic("not implemented")
}
