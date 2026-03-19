package ports

import (
	"context"
	"roomService/internal/core/domain"
)

type RoomService interface {
	CreateRoom(ctx context.Context, name string, isPrivate bool, groupID string) (domain.CreateRoomResult, error)
	GetRoom(ctx context.Context, roomUUID string) (domain.Room, error)
	UpdateRoomName(ctx context.Context, roomUUID string, name string) error
	DeleteRoom(ctx context.Context, roomUUID string) error
	GetUserRooms(ctx context.Context) ([]domain.UserRoom, error)

	GetRoomMembers(ctx context.Context, roomUUID string) ([]domain.RoomMember, error)
	LeaveRoom(ctx context.Context, roomUUID string) error
	RemoveMember(ctx context.Context, roomUUID string, userUUID string) error

	CreateInvite(ctx context.Context, roomUUID string, usesLeft int32, expiresAt int64) (domain.CreateInviteResult, error)
	GetInvite(ctx context.Context, inviteID string) (domain.RoomInvite, error)
	DeleteInvite(ctx context.Context, inviteID string, roomUUID string) error
	JoinViaInvite(ctx context.Context, inviteID string) error

	CreateJoinRequest(ctx context.Context, roomUUID string, publicKey string) error
	GetJoinRequests(ctx context.Context, roomUUID string) ([]domain.JoinRequest, error)
	RespondToJoinRequest(ctx context.Context, roomUUID string, userUUID string, decision domain.RequestStatus) error
	GetPendingKeys(ctx context.Context) ([]domain.PendingKey, error)
	AcknowledgeKeyDelivery(ctx context.Context, roomUUID string) error

	MarkAsRead(ctx context.Context, roomUUID string, lastReadMessageID string) error
	GetUnreadCount(ctx context.Context, roomUUID string) (domain.GetUnreadCountResult, error)

	GetAllowedRooms(ctx context.Context, userUUID string) ([]string, error)
}
