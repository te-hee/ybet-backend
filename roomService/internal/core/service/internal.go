package service

import (
	"context"
	"roomService/internal/contextkeys"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *roomService) GetAllowedRooms(ctx context.Context, userUUID string) ([]string, error) {
	var roomsIDs []string
	if !contextkeys.IsInternalRequest(ctx) {
		return roomsIDs, status.Errorf(codes.Unauthenticated, "internal requests only")
	}
	rooms, err := s.repo.GetUserRooms(ctx, userUUID)
	if err != nil {
		return roomsIDs, err
	}

	for _, room := range rooms {
		roomsIDs = append(roomsIDs, room.RoomUUID)
	}

	return roomsIDs, nil
}
