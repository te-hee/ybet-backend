package service

import (
	"gateway/internal/client"
	"gateway/internal/model"
)

type RoomService struct {
	client client.RoomClient
}

func NewRoomService(client client.RoomClient) *RoomService {
	return &RoomService{
		client: client,
	}
}

// ─── Room CRUD ──────────────────────────────────────────────────────

func (s *RoomService) CreateRoom(userToken string, req model.CreateRoomRequest) (*model.CreateRoomResponse, error) {
	return s.client.CreateRoom(userToken, req)
}

func (s *RoomService) GetRoom(userToken string, roomUUID string) (*model.RoomResponse, error) {
	return s.client.GetRoom(userToken, roomUUID)
}

func (s *RoomService) UpdateRoomName(userToken string, req model.UpdateRoomNameRequest) error {
	return s.client.UpdateRoomName(userToken, req)
}

func (s *RoomService) DeleteRoom(userToken string, roomUUID string) error {
	return s.client.DeleteRoom(userToken, roomUUID)
}

func (s *RoomService) GetUserRooms(userToken string) ([]model.UserRoom, error) {
	return s.client.GetUserRooms(userToken)
}

// ─── Membership ─────────────────────────────────────────────────────

func (s *RoomService) GetRoomMembers(userToken string, roomUUID string) ([]model.RoomMember, error) {
	return s.client.GetRoomMembers(userToken, roomUUID)
}

func (s *RoomService) LeaveRoom(userToken string, roomUUID string) error {
	return s.client.LeaveRoom(userToken, roomUUID)
}

func (s *RoomService) RemoveMember(userToken string, req model.RemoveMemberRequest) error {
	return s.client.RemoveMember(userToken, req)
}

// ─── Invite Links ───────────────────────────────────────────────────

func (s *RoomService) CreateInvite(userToken string, req model.CreateInviteRequest) (*model.CreateInviteResponse, error) {
	return s.client.CreateInvite(userToken, req)
}

func (s *RoomService) GetInvite(userToken string, inviteID string) (*model.InviteResponse, error) {
	return s.client.GetInvite(userToken, inviteID)
}

func (s *RoomService) DeleteInvite(userToken string, req model.DeleteInviteRequest) error {
	return s.client.DeleteInvite(userToken, req)
}

func (s *RoomService) JoinViaInvite(userToken string, inviteID string) error {
	return s.client.JoinViaInvite(userToken, inviteID)
}

// ─── Join Requests ──────────────────────────────────────────────────

func (s *RoomService) CreateJoinRequest(userToken string, req model.CreateJoinRequestRequest) error {
	return s.client.CreateJoinRequest(userToken, req)
}

func (s *RoomService) GetJoinRequests(userToken string, roomUUID string) ([]model.JoinRequestResponse, error) {
	return s.client.GetJoinRequests(userToken, roomUUID)
}

func (s *RoomService) RespondToJoinRequest(userToken string, req model.RespondToJoinRequestRequest) error {
	return s.client.RespondToJoinRequest(userToken, req)
}

// ─── Unread Tracking ────────────────────────────────────────────────

func (s *RoomService) MarkAsRead(userToken string, req model.MarkAsReadRequest) error {
	return s.client.MarkAsRead(userToken, req)
}

func (s *RoomService) GetUnreadCount(userToken string, roomUUID string) (*model.UnreadCountResponse, error) {
	return s.client.GetUnreadCount(userToken, roomUUID)
}
