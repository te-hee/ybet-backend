package model

// ─── Room CRUD ──────────────────────────────────────────────────────

type CreateRoomRequest struct {
	Name      string `json:"name"`
	IsPrivate bool   `json:"is_private"`
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
	RoomUUID string `json:"room_uuid"`
	Name     string `json:"name"`
}

type UserRoom struct {
	RoomUUID    string `json:"room_uuid"`
	Name        string `json:"name"`
	IsPrivate   bool   `json:"is_private"`
	UnreadCount int32  `json:"unread_count"`
	JoinedAt    int64  `json:"joined_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// ─── Membership ─────────────────────────────────────────────────────

type RoomMember struct {
	UserUUID string `json:"user_uuid"`
	JoinedAt int64  `json:"joined_at"`
}

type RemoveMemberRequest struct {
	RoomUUID string `json:"room_uuid"`
	UserUUID string `json:"user_uuid"`
}

// ─── Invite Links ───────────────────────────────────────────────────

type CreateInviteRequest struct {
	RoomUUID  string `json:"room_uuid"`
	UsesLeft  int32  `json:"uses_left"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
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
	InviteID string `json:"invite_id"`
	RoomUUID string `json:"room_uuid"`
}

// ─── Join Requests ──────────────────────────────────────────────────

type CreateJoinRequestRequest struct {
	RoomUUID  string `json:"room_uuid"`
	PublicKey string `json:"public_key"`
}

type JoinRequestResponse struct {
	RoomUUID  string `json:"room_uuid"`
	UserUUID  string `json:"user_uuid"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type RespondToJoinRequestRequest struct {
	RoomUUID string `json:"room_uuid"`
	UserUUID string `json:"user_uuid"`
	Decision string `json:"decision"`
}

// ─── Unread Tracking ────────────────────────────────────────────────

type MarkAsReadRequest struct {
	RoomUUID          string `json:"room_uuid"`
	LastReadMessageID string `json:"last_read_message_id"`
}

type UnreadCountResponse struct {
	UnreadCount int32 `json:"unread_count"`
}
