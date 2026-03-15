package domain

import "time"

type RequestStatus int32

const (
	RequestStatusUnspecified RequestStatus = 0
	RequestStatusPending     RequestStatus = 1
	RequestStatusAccepted    RequestStatus = 2
	RequestStatusRejected    RequestStatus = 3
)

type Room struct {
	RoomUUID    string
	Name        string
	AdminID     string
	IsPrivate   bool
	GroupID     string
	MemberCount int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoomMember struct {
	UserUUID string
	JoinedAt time.Time
}

type RoomInvite struct {
	InviteID  string
	RoomUUID  string
	UsesLeft  int32
	ExpiresAt time.Time
	CreatedAt time.Time
}

type JoinRequest struct {
	RoomUUID  string
	UserUUID  string
	Username  string
	PublicKey string
	Status    RequestStatus
	CreatedAt time.Time
}

type UserRoom struct {
	RoomUUID    string
	Name        string
	IsPrivate   bool
	UnreadCount int32
	JoinedAt    time.Time
	UpdatedAt   time.Time
}

type CreateRoomResult struct {
	RoomUUID  string
	CreatedAt time.Time
}

type CreateInviteResult struct {
	InviteID string
}

type GetUnreadCountResult struct {
	UnreadCount int32
}
