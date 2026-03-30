package model

// ─── Room CRUD ──────────────────────────────────────────────────────

type CreateRoomRequest struct {
	Name      string `json:"name" validate:"required"`
	IsPrivate bool   `json:"is_private" validate:"required"`
	GroupId   string `json:"group_id,omitempty"`
}

type CreateRoomResponse struct {
	RoomUUID  string `json:"room_uuid"`
	CreatedAt int64  `json:"created_at"`
}

type RoomResponse struct {
	RoomUUID    string `json:"room_uuid"`
	Name        string `json:"name"`
	AdminID     string `json:"admin_id"`
	IsPrivate   bool   `json:"is_private"`
	GroupID     string `json:"group_id,omitempty"`
	MemberCount int32  `json:"member_count"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type UpdateRoomNameRequest struct {
	RoomUUID string `json:"room_uuid" validate:"required,uuid"`
	Name     string `json:"name" validate:"required"`
}

type UserRoom struct {
	RoomUUID    string `json:"room_uuid"`
	Name        string `json:"name"`
	IsPrivate   bool   `json:"is_private"`
	UnreadCount int32  `json:"unread_count"`
	JoinedAt    int64  `json:"joined_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type GetRoomRequest struct{
	RoomUUID string `query:"room_uuid" validate:"required,uuid"`
}

type DeleteRoomRequst struct{
	RoomUUID string `query:"room_uuid" validate:"required,uuid"`
}


// ─── Membership ─────────────────────────────────────────────────────

type RoomMember struct {
	UserUUID string `json:"user_uuid"`
	JoinedAt int64  `json:"joined_at"`
}

type RemoveMemberRequest struct {
	RoomUUID string `json:"room_uuid" validate:"required,uuid"`
	UserUUID string `json:"user_uuid" validate:"required,uuid"`
}

type GetRoomMembersRequest struct{
	RoomUUID string `query:"room_uuid" validate:"required,uuid"`
}

type LeaveRoomRequest struct{
	RoomUUID string `json:"room_uuid" validate:"required,uuid"`
}

// ─── Invite Links ───────────────────────────────────────────────────

type CreateInviteRequest struct {
	RoomUUID  string `json:"room_uuid" validate:"required,uuid"`
	UsesLeft  int32  `json:"uses_left" validate:"required"`
	ExpiresAt int64  `json:"expires_at,omitempty" validate:"required"`
}

type CreateInviteResponse struct {
	InviteID string `json:"invite_id"`
}

type InviteResponse struct {
	InviteID  string `json:"invite_id"`
	RoomUUID  string `json:"room_uuid"`
	UsesLeft  int32  `json:"uses_left"`
	ExpiresAt int64  `json:"expires_at"`
	CreatedAt int64  `json:"created_at"`
}

type DeleteInviteRequest struct {
	InviteID string `json:"invite_id" validate:"required,uuid"`
	RoomUUID string `json:"room_uuid" validate:"required,uuid"`
}


type GetInviteRequest struct{
	InvieID string `query:"invite_id" validate:"required,uuid"`
}

type JoinViaInviteRequest struct{
	InviteID string `json:"invite_id" validate:"required,uuid"`
}

// ─── Join Requests ──────────────────────────────────────────────────

type CreateJoinRequestRequest struct {
	RoomUUID  string `json:"room_uuid" validate:"required,uuid"`
	PublicKey string `json:"public_key" validate:"required"`
}

type JoinRequestResponse struct {
	RoomUUID  string `json:"room_uuid"`
	UserUUID  string `json:"user_uuid"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type Decision string

const(
	ACCEPTED Decision = "ACCEPTED"
	REJECTED Decision = "REJECTED"
)

type RespondToJoinRequestRequest struct {
	RoomUUID string `json:"room_uuid" validate:"required,uuid"`
	UserUUID string `json:"user_uuid" validate:"required,uuid"`
	Decision Decision `json:"decision" validate:"required,oneof=ACCEPTED REJECTED"`
}

type GetJoinReqeustRequest struct{
	RoomUUID string `query:"room_uuid" validate:"required,uuid"`
}


// ─── Unread Tracking ────────────────────────────────────────────────

type MarkAsReadRequest struct {
	RoomUUID          string `json:"room_uuid" validate:"required,uuid"`
	LastReadMessageID string `json:"last_read_message_id" validate:"required,uuid"`
}

type UnreadCountResponse struct {
	UnreadCount int32 `json:"unread_count"`
}

type UnreadCountReqeust struct{
	RoomUUID string `query:"room_uuid" validate:"required,uuid"`
}
