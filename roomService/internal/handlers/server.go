package handlers

import (
	roomv1 "backend/proto/room/v1"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RoomServer struct {
	roomv1.UnimplementedRoomServiceServer
}

func NewRoomServer() *RoomServer {
	return &RoomServer{}
}

func (s *RoomServer) CreateRoom(_ context.Context, _ *roomv1.CreateRoomRequest) (*roomv1.CreateRoomResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateRoom not implemented")
}

func (s *RoomServer) GetRoom(_ context.Context, _ *roomv1.GetRoomRequest) (*roomv1.GetRoomResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetRoom not implemented")
}

func (s *RoomServer) UpdateRoomName(_ context.Context, _ *roomv1.UpdateRoomNameRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "UpdateRoomName not implemented")
}

func (s *RoomServer) DeleteRoom(_ context.Context, _ *roomv1.DeleteRoomRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteRoom not implemented")
}

func (s *RoomServer) GetUserRooms(_ context.Context, _ *roomv1.GetUserRoomsRequest) (*roomv1.GetUserRoomsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetUserRooms not implemented")
}

func (s *RoomServer) GetRoomMembers(_ context.Context, _ *roomv1.GetRoomMembersRequest) (*roomv1.GetRoomMembersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetRoomMembers not implemented")
}

func (s *RoomServer) LeaveRoom(_ context.Context, _ *roomv1.LeaveRoomRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "LeaveRoom not implemented")
}

func (s *RoomServer) RemoveMember(_ context.Context, _ *roomv1.RemoveMemberRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "RemoveMember not implemented")
}

func (s *RoomServer) CreateInvite(_ context.Context, _ *roomv1.CreateInviteRequest) (*roomv1.CreateInviteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateInvite not implemented")
}

func (s *RoomServer) GetInvite(_ context.Context, _ *roomv1.GetInviteRequest) (*roomv1.GetInviteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetInvite not implemented")
}

func (s *RoomServer) DeleteInvite(_ context.Context, _ *roomv1.DeleteInviteRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteInvite not implemented")
}

func (s *RoomServer) JoinViaInvite(_ context.Context, _ *roomv1.JoinViaInviteRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "JoinViaInvite not implemented")
}

func (s *RoomServer) CreateJoinRequest(_ context.Context, _ *roomv1.CreateJoinRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "CreateJoinRequest not implemented")
}

func (s *RoomServer) GetJoinRequests(_ context.Context, _ *roomv1.GetJoinRequestsRequest) (*roomv1.GetJoinRequestsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetJoinRequests not implemented")
}

func (s *RoomServer) RespondToJoinRequest(_ context.Context, _ *roomv1.RespondToJoinRequestRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "RespondToJoinRequest not implemented")
}

func (s *RoomServer) MarkAsRead(_ context.Context, _ *roomv1.MarkAsReadRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "MarkAsRead not implemented")
}

func (s *RoomServer) GetUnreadCount(_ context.Context, _ *roomv1.GetUnreadCountRequest) (*roomv1.GetUnreadCountResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetUnreadCount not implemented")
}

func (s *RoomServer) GetAllowedRooms(_ context.Context, _ *roomv1.GetAllowedRoomsRequest) (*roomv1.GetAllowedRoomsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetAllowedRooms not implemented")
}
