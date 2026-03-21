package service

import (
	"context"
	"roomService/internal/contextkeys"
	"roomService/internal/core/domain"
)

func (s *roomService) GetAllowedRooms(ctx context.Context, req domain.GetAllowedRoomsDTO) ([]string, error) {
	var roomsIDs []string
	if !contextkeys.IsInternalRequest(ctx) {
		return roomsIDs, domain.NewError(domain.CodeUnauthenticated, "internal requests only")
	}
	rooms, err := s.repo.GetUserRooms(ctx, req.UserUUID)
	if err != nil {
		return roomsIDs, err
	}

	for _, room := range rooms {
		roomsIDs = append(roomsIDs, room.RoomUUID)
	}

	return roomsIDs, nil
}
