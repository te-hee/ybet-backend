package domain

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateRoomDTO struct {
	Name      string `validate:"required,min=1,max=100"`
	IsPrivate bool
	GroupID   string `validate:"omitempty,uuid"`
}

type GetRoomDTO struct {
	RoomUUID string `validate:"required,uuid"`
}

type UpdateRoomNameDTO struct {
	RoomUUID string `validate:"required,uuid"`
	Name     string `validate:"required,min=1,max=100"`
}

type DeleteRoomDTO struct {
	RoomUUID string `validate:"required,uuid"`
}

type GetRoomMembersDTO struct {
	RoomUUID string `validate:"required,uuid"`
}

type LeaveRoomDTO struct {
	RoomUUID string `validate:"required,uuid"`
}

type RemoveMemberDTO struct {
	RoomUUID string `validate:"required,uuid"`
	UserUUID string `validate:"required,uuid"`
}

type CreateInviteDTO struct {
	RoomUUID  string `validate:"required,uuid"`
	UsesLeft  int32
	ExpiresAt *timestamppb.Timestamp
}

type GetInviteDTO struct {
	InviteID string `validate:"required"`
}

type DeleteInviteDTO struct {
	InviteID string `validate:"required"`
	RoomUUID string `validate:"required,uuid"`
}

type JoinViaInviteDTO struct {
	InviteID string `validate:"required"`
}

type CreateJoinRequestDTO struct {
	RoomUUID  string `validate:"required,uuid"`
	PublicKey string `validate:"required"`
}

type GetJoinRequestsDTO struct {
	RoomUUID string `validate:"required,uuid"`
}

type RespondToJoinRequestDTO struct {
	RoomUUID     string        `validate:"required,uuid"`
	UserUUID     string        `validate:"required,uuid"`
	Decision     RequestStatus `validate:"required"`
	EncryptedKey string        `validate:"required"`
}

type AcknowledgeKeyDeliveryDTO struct {
	RoomUUID string `validate:"required,uuid"`
}

type MarkAsReadDTO struct {
	RoomUUID          string `validate:"required,uuid"`
	LastReadMessageID string `validate:"required"`
}

type GetUnreadCountDTO struct {
	RoomUUID string `validate:"required,uuid"`
}

type GetAllowedRoomsDTO struct {
	UserUUID string `validate:"required,uuid"`
}
