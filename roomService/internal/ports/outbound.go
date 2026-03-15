package ports

import (
	"context"
	"roomService/internal/core/domain"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, room domain.Room) error
	GetRoom(ctx context.Context, roomUUID string) (domain.Room, error)
	UpdateRoomName(ctx context.Context, roomUUID string, name string) error
	DeleteRoom(ctx context.Context, roomUUID string) error
	GetUserRooms(ctx context.Context, userUUID string) ([]domain.UserRoom, error)

	GetRoomMembers(ctx context.Context, roomUUID string) ([]domain.RoomMember, error)
	AddMember(ctx context.Context, roomUUID string, userUUID string) error
	RemoveMember(ctx context.Context, roomUUID string, userUUID string) error

	CreateInvite(ctx context.Context, invite domain.RoomInvite) error
	GetInvite(ctx context.Context, inviteID string) (domain.RoomInvite, error)
	DeleteInvite(ctx context.Context, inviteID string, roomUUID string) error
	DecrementInviteUses(ctx context.Context, inviteID string) error

	CreateJoinRequest(ctx context.Context, req domain.JoinRequest) error
	GetJoinRequests(ctx context.Context, roomUUID string) ([]domain.JoinRequest, error)
	UpdateJoinRequestStatus(ctx context.Context, roomUUID string, userUUID string, status domain.RequestStatus) error

	MarkAsRead(ctx context.Context, roomUUID string, userUUID string, lastReadMessageID string) error
	GetUnreadCount(ctx context.Context, roomUUID string, userUUID string) (int32, error)
}

type EventPublisher interface {
	PublishRoomCreated(ctx context.Context, room domain.Room) error
	PublishRoomDeleted(ctx context.Context, roomUUID string) error
	PublishMemberJoined(ctx context.Context, roomUUID string, userUUID string) error
	PublishMemberLeft(ctx context.Context, roomUUID string, userUUID string) error
}
