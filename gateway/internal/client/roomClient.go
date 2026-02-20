package client

import (
	roomv1 "backend/proto/room/v1"
	"context"
	"gateway/config"

	"google.golang.org/grpc/metadata"
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

func (c *RoomServiceClient) CreateRoom(req *roomv1.CreateRoomRequest) (*roomv1.CreateRoomResponse, error) {
	ctxWithAuth := setRoomServiceAuth(c.ctx)
	return c.grpcClient.CreateRoom(ctxWithAuth, req)
}

func (c *RoomServiceClient) GetRoom(req *roomv1.GetRoomRequest) (*roomv1.GetRoomResponse, error) {
	ctxWithAuth := setRoomServiceAuth(c.ctx)
	return c.grpcClient.GetRoom(ctxWithAuth, req)
}

func (c *RoomServiceClient) UpdateRoomName(req *roomv1.UpdateRoomNameRequest) error {
	ctxWithAuth := setRoomServiceAuth(c.ctx)
	_, err := c.grpcClient.UpdateRoomName(ctxWithAuth, req)
	return err
}

func (c *RoomServiceClient) DeleteRoom(req *roomv1.DeleteRoomRequest) error {
	ctxWithAuth := setRoomServiceAuth(c.ctx)
	_, err := c.grpcClient.DeleteRoom(ctxWithAuth, req)
	return err
}

func (c *RoomServiceClient) JoinRoom(req *roomv1.JoinRoomRequest) error {
	ctxWithAuth := setRoomServiceAuth(c.ctx)
	_, err := c.grpcClient.JoinRoom(ctxWithAuth, req)
	return err
}

func (c *RoomServiceClient) LeaveRoom(req *roomv1.LeaveRoomRequest) error {
	ctxWithAuth := setRoomServiceAuth(c.ctx)
	_, err := c.grpcClient.LeaveRoom(ctxWithAuth, req)
	return err
}

func (c *RoomServiceClient) GetRoomMembers(req *roomv1.GetRoomMembersRequest) (*roomv1.GetRoomMembersResponse, error) {
	ctxWithAuth := setRoomServiceAuth(c.ctx)
	return c.grpcClient.GetRoomMembers(ctxWithAuth, req)
}

func (c *RoomServiceClient) GetUserRooms(req *roomv1.GetUserRoomsRequest) (*roomv1.GetUserRoomsResponse, error) {
	ctxWithAuth := setRoomServiceAuth(c.ctx)
	return c.grpcClient.GetUserRooms(ctxWithAuth, req)
}

func (c *RoomServiceClient) MarkAsRead(req *roomv1.MarkAsReadRequest) error {
	ctxWithAuth := setRoomServiceAuth(c.ctx)
	_, err := c.grpcClient.MarkAsRead(ctxWithAuth, req)
	return err
}

func (c *RoomServiceClient) GetUnreadCount(req *roomv1.GetUnreadCountRequest) (*roomv1.GetUnreadCountResponse, error) {
	ctxWithAuth := setRoomServiceAuth(c.ctx)
	return c.grpcClient.GetUnreadCount(ctxWithAuth, req)
}

func setRoomServiceAuth(ctx context.Context) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+*config.RoomServiceKey)
	ctxWithAuth := metadata.NewOutgoingContext(ctx, md)
	return ctxWithAuth
}
