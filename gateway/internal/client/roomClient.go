package client

import (
	roomv1 "backend/proto/room/v1"
	"context"
	"gateway/config"
	"gateway/internal/model"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RoomServiceClient struct {
	ctx        context.Context
	grpcClient roomv1.RoomServiceClient
}

func NewRoomServiceClient(ctx context.Context, grpcClient roomv1.RoomServiceClient) *RoomServiceClient {
	return &RoomServiceClient{
		ctx:        ctx,
		grpcClient: grpcClient,
	}
}

// ─── Room CRUD ──────────────────────────────────────────────────────

func (c *RoomServiceClient) CreateRoom(userToken string, req model.CreateRoomRequest) (*model.CreateRoomResponse, error) {
	ctx := setRoomAuth(c.ctx, userToken)
	resp, err := c.grpcClient.CreateRoom(ctx, &roomv1.CreateRoomRequest{
		Name:      req.Name,
		IsPrivate: req.IsPrivate,
		GroupId:   req.GroupId,
	})
	if err != nil {
		return nil, err
	}
	return &model.CreateRoomResponse{
		RoomUUID:  resp.RoomUuid,
		CreatedAt: resp.CreatedAt.AsTime().Unix(),
	}, nil
}

func (c *RoomServiceClient) GetRoom(userToken string, roomUUID string) (*model.RoomResponse, error) {
	ctx := setRoomAuth(c.ctx, userToken)
	resp, err := c.grpcClient.GetRoom(ctx, &roomv1.GetRoomRequest{RoomUuid: roomUUID})
	if err != nil {
		return nil, err
	}
	return protoRoomToModel(resp.Room), nil
}

func (c *RoomServiceClient) UpdateRoomName(userToken string, req model.UpdateRoomNameRequest) error {
	ctx := setRoomAuth(c.ctx, userToken)
	_, err := c.grpcClient.UpdateRoomName(ctx, &roomv1.UpdateRoomNameRequest{
		RoomUuid: req.RoomUUID,
		Name:     req.Name,
	})
	return err
}

func (c *RoomServiceClient) DeleteRoom(userToken string, roomUUID string) error {
	ctx := setRoomAuth(c.ctx, userToken)
	_, err := c.grpcClient.DeleteRoom(ctx, &roomv1.DeleteRoomRequest{RoomUuid: roomUUID})
	return err
}

func (c *RoomServiceClient) GetUserRooms(userToken string) ([]model.UserRoom, error) {
	ctx := setRoomAuth(c.ctx, userToken)
	resp, err := c.grpcClient.GetUserRooms(ctx, &roomv1.GetUserRoomsRequest{})
	if err != nil {
		return nil, err
	}
	rooms := make([]model.UserRoom, 0, len(resp.Rooms))
	for _, r := range resp.Rooms {
		rooms = append(rooms, model.UserRoom{
			RoomUUID:    r.RoomUuid,
			Name:        r.Name,
			IsPrivate:   r.IsPrivate,
			UnreadCount: r.UnreadCount,
			JoinedAt:    r.JoinedAt.AsTime().Unix(),
			UpdatedAt:   r.UpdatedAt.AsTime().Unix(),
		})
	}
	return rooms, nil
}

// ─── Membership ─────────────────────────────────────────────────────

func (c *RoomServiceClient) GetRoomMembers(userToken string, roomUUID string) ([]model.RoomMember, error) {
	ctx := setRoomAuth(c.ctx, userToken)
	resp, err := c.grpcClient.GetRoomMembers(ctx, &roomv1.GetRoomMembersRequest{RoomUuid: roomUUID})
	if err != nil {
		return nil, err
	}
	members := make([]model.RoomMember, 0, len(resp.Members))
	for _, m := range resp.Members {
		members = append(members, model.RoomMember{
			UserUUID: m.UserUuid,
			JoinedAt: m.JoinedAt.AsTime().Unix(),
		})
	}
	return members, nil
}

func (c *RoomServiceClient) LeaveRoom(userToken string, roomUUID string) error {
	ctx := setRoomAuth(c.ctx, userToken)
	_, err := c.grpcClient.LeaveRoom(ctx, &roomv1.LeaveRoomRequest{RoomUuid: roomUUID})
	return err
}

func (c *RoomServiceClient) RemoveMember(userToken string, req model.RemoveMemberRequest) error {
	ctx := setRoomAuth(c.ctx, userToken)
	_, err := c.grpcClient.RemoveMember(ctx, &roomv1.RemoveMemberRequest{
		RoomUuid: req.RoomUUID,
		UserUuid: req.UserUUID,
	})
	return err
}

// ─── Invite Links ───────────────────────────────────────────────────

func (c *RoomServiceClient) CreateInvite(userToken string, req model.CreateInviteRequest) (*model.CreateInviteResponse, error) {
	ctx := setRoomAuth(c.ctx, userToken)
	protoReq := &roomv1.CreateInviteRequest{
		RoomUuid: req.RoomUUID,
		UsesLeft: req.UsesLeft,
	}
	if req.ExpiresAt != 0 {
		protoReq.ExpiresAt = timestamppb.New(timestampFromUnix(req.ExpiresAt))
	}
	resp, err := c.grpcClient.CreateInvite(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return &model.CreateInviteResponse{InviteID: resp.InviteId}, nil
}

func (c *RoomServiceClient) GetInvite(userToken string, inviteID string) (*model.InviteResponse, error) {
	ctx := setRoomAuth(c.ctx, userToken)
	resp, err := c.grpcClient.GetInvite(ctx, &roomv1.GetInviteRequest{InviteId: inviteID})
	if err != nil {
		return nil, err
	}
	inv := resp.Invite
	return &model.InviteResponse{
		InviteID:  inv.InviteId,
		RoomUUID:  inv.RoomUuid,
		UsesLeft:  inv.UsesLeft,
		ExpiresAt: inv.ExpiresAt.AsTime().Unix(),
		CreatedAt: inv.CreatedAt.AsTime().Unix(),
	}, nil
}

func (c *RoomServiceClient) DeleteInvite(userToken string, req model.DeleteInviteRequest) error {
	ctx := setRoomAuth(c.ctx, userToken)
	_, err := c.grpcClient.DeleteInvite(ctx, &roomv1.DeleteInviteRequest{
		InviteId: req.InviteID,
		RoomUuid: req.RoomUUID,
	})
	return err
}

func (c *RoomServiceClient) JoinViaInvite(userToken string, inviteID string) error {
	ctx := setRoomAuth(c.ctx, userToken)
	_, err := c.grpcClient.JoinViaInvite(ctx, &roomv1.JoinViaInviteRequest{InviteId: inviteID})
	return err
}

// ─── Join Requests ──────────────────────────────────────────────────

func (c *RoomServiceClient) CreateJoinRequest(userToken string, req model.CreateJoinRequestRequest) error {
	ctx := setRoomAuth(c.ctx, userToken)
	_, err := c.grpcClient.CreateJoinRequest(ctx, &roomv1.CreateJoinRequestRequest{
		RoomUuid:  req.RoomUUID,
		PublicKey: req.PublicKey,
	})
	return err
}

func (c *RoomServiceClient) GetJoinRequests(userToken string, roomUUID string) ([]model.JoinRequestResponse, error) {
	ctx := setRoomAuth(c.ctx, userToken)
	resp, err := c.grpcClient.GetJoinRequests(ctx, &roomv1.GetJoinRequestsRequest{RoomUuid: roomUUID})
	if err != nil {
		return nil, err
	}
	requests := make([]model.JoinRequestResponse, 0, len(resp.Requests))
	for _, r := range resp.Requests {
		requests = append(requests, model.JoinRequestResponse{
			RoomUUID:  r.RoomUuid,
			UserUUID:  r.UserUuid,
			Username:  r.Username,
			PublicKey: r.PublicKey,
			Status:    r.Status.String(),
			CreatedAt: r.CreatedAt.AsTime().Unix(),
		})
	}
	return requests, nil
}

func (c *RoomServiceClient) RespondToJoinRequest(userToken string, req model.RespondToJoinRequestRequest) error {
	ctx := setRoomAuth(c.ctx, userToken)
	decision := roomv1.RequestStatus_REQUEST_STATUS_REJECTED
	if req.Decision == "ACCEPTED" {
		decision = roomv1.RequestStatus_REQUEST_STATUS_ACCEPTED
	}
	_, err := c.grpcClient.RespondToJoinRequest(ctx, &roomv1.RespondToJoinRequestRequest{
		RoomUuid: req.RoomUUID,
		UserUuid: req.UserUUID,
		Decision: decision,
	})
	return err
}

// ─── Unread Tracking ────────────────────────────────────────────────

func (c *RoomServiceClient) MarkAsRead(userToken string, req model.MarkAsReadRequest) error {
	ctx := setRoomAuth(c.ctx, userToken)
	_, err := c.grpcClient.MarkAsRead(ctx, &roomv1.MarkAsReadRequest{
		RoomUuid:          req.RoomUUID,
		LastReadMessageId: req.LastReadMessageID,
	})
	return err
}

func (c *RoomServiceClient) GetUnreadCount(userToken string, roomUUID string) (*model.UnreadCountResponse, error) {
	ctx := setRoomAuth(c.ctx, userToken)
	resp, err := c.grpcClient.GetUnreadCount(ctx, &roomv1.GetUnreadCountRequest{RoomUuid: roomUUID})
	if err != nil {
		return nil, err
	}
	return &model.UnreadCountResponse{UnreadCount: resp.UnreadCount}, nil
}

// ─── Internal ───────────────────────────────────────────────────────

func (c *RoomServiceClient) GetAllowedRooms(userUUID string) ([]string, error) {
	ctx := setRoomInternalAuth(c.ctx)
	resp, err := c.grpcClient.GetAllowedRooms(ctx, &roomv1.GetAllowedRoomsRequest{UserUuid: userUUID})
	if err != nil {
		return nil, err
	}
	return resp.RoomUuids, nil
}

// ─── Helpers ────────────────────────────────────────────────────────

// setRoomInternalAuth sets only the internal API key (no JWT) — for service-to-service RPCs.
func setRoomInternalAuth(ctx context.Context) context.Context {
	md := metadata.Pairs("x-internal-api-key", config.Cfg.Services.Room.ApiKey)
	return metadata.NewOutgoingContext(ctx, md)
}

// setRoomAuth sets the internal API key and forwards the user's JWT — for user-facing RPCs.
func setRoomAuth(ctx context.Context, userToken string) context.Context {
	md := metadata.Pairs(
		"x-internal-api-key", config.Cfg.Services.Room.ApiKey,
		"authorization", "Bearer "+userToken,
	)
	return metadata.NewOutgoingContext(ctx, md)
}

func protoRoomToModel(r *roomv1.Room) *model.RoomResponse {
	return &model.RoomResponse{
		RoomUUID:    r.RoomUuid,
		Name:        r.Name,
		AdminID:     r.AdminId,
		IsPrivate:   r.IsPrivate,
		GroupID:     r.GroupId,
		MemberCount: r.MemberCount,
		CreatedAt:   r.CreatedAt.AsTime().Unix(),
		UpdatedAt:   r.UpdatedAt.AsTime().Unix(),
	}
}

func timestampFromUnix(unix int64) time.Time {
	return time.Unix(unix, 0)
}
