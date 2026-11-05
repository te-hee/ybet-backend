package ports

import (
	"context"
	"roomService/internal/core/domain"
)

type RoomService interface {
	CreateRoom(ctx context.Context, req domain.CreateRoomDTO) (domain.CreateRoomResult, error)
	GetRoom(ctx context.Context, req domain.GetRoomDTO) (domain.Room, error)
	UpdateRoomName(ctx context.Context, req domain.UpdateRoomNameDTO) error
	DeleteRoom(ctx context.Context, req domain.DeleteRoomDTO) error
	GetUserRooms(ctx context.Context) ([]domain.UserRoom, error)

	GetRoomMembers(ctx context.Context, req domain.GetRoomMembersDTO) ([]domain.RoomMember, error)
	LeaveRoom(ctx context.Context, req domain.LeaveRoomDTO) error
	RemoveMember(ctx context.Context, req domain.RemoveMemberDTO) error

	CreateInvite(ctx context.Context, req domain.CreateInviteDTO) (domain.CreateInviteResult, error)
	GetInvite(ctx context.Context, req domain.GetInviteDTO) (domain.RoomInvite, error)
	DeleteInvite(ctx context.Context, req domain.DeleteInviteDTO) error
	JoinViaInvite(ctx context.Context, req domain.JoinViaInviteDTO) error

	CreateJoinRequest(ctx context.Context, req domain.CreateJoinRequestDTO) error
	GetJoinRequests(ctx context.Context, req domain.GetJoinRequestsDTO) ([]domain.JoinRequest, error)
	RespondToJoinRequest(ctx context.Context, req domain.RespondToJoinRequestDTO) error
	GetPendingKeys(ctx context.Context) ([]domain.PendingKey, error)
	AcknowledgeKeyDelivery(ctx context.Context, req domain.AcknowledgeKeyDeliveryDTO) error

	MarkAsRead(ctx context.Context, req domain.MarkAsReadDTO) error
	GetUnreadCount(ctx context.Context, req domain.GetUnreadCountDTO) (domain.GetUnreadCountResult, error)

	GetAllowedRooms(ctx context.Context, req domain.GetAllowedRoomsDTO) ([]string, error)
}
