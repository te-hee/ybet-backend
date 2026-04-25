package client

import (
	"gateway/internal/model"
)

type MessageClient interface {
	GetMessageHistory(limit uint32) ([]model.Message, error)
	SendMessage(message model.SendMessageRequest) (*model.OutputSendMessege, error)
	EditMessage(editRequest model.EditMessageRequest) error
	DeleteMessage(deleteRequest model.DeleteMessageRequest) error
}

type RoomClient interface {
	// Room CRUD
	CreateRoom(userToken string, req model.CreateRoomRequest) (*model.CreateRoomResponse, error)
	GetRoom(userToken string, roomUUID string) (*model.RoomResponse, error)
	UpdateRoomName(userToken string, req model.UpdateRoomNameRequest) error
	DeleteRoom(userToken string, roomUUID string) error
	GetUserRooms(userToken string) ([]model.UserRoom, error)

	// Membership
	GetRoomMembers(userToken string, roomUUID string) ([]model.RoomMember, error)
	LeaveRoom(userToken string, roomUUID string) error
	RemoveMember(userToken string, req model.RemoveMemberRequest) error

	// Invite Links
	CreateInvite(userToken string, req model.CreateInviteRequest) (*model.CreateInviteResponse, error)
	GetInvite(userToken string, inviteID string) (*model.InviteResponse, error)
	DeleteInvite(userToken string, req model.DeleteInviteRequest) error
	JoinViaInvite(userToken string, inviteID string) error

	// Join Requests
	CreateJoinRequest(userToken string, req model.CreateJoinRequestRequest) error
	GetJoinRequests(userToken string, roomUUID string) ([]model.JoinRequestResponse, error)
	RespondToJoinRequest(userToken string, req model.RespondToJoinRequestRequest) error

	// Unread Tracking
	MarkAsRead(userToken string, req model.MarkAsReadRequest) error
	GetUnreadCount(userToken string, roomUUID string) (*model.UnreadCountResponse, error)

	// Internal
	GetAllowedRooms(userUUID string) ([]string, error)
}

type UserClient interface {
	SignIn(password string, username string) (*model.SignInResponseV2, error)
	LogIn(password string, username string) (*model.LogInResponseV2, error)
	GetNewAuthToken(refreshToken string) (*string, error)
}
