package grpcadapter

import (
	roomv1 "backend/proto/room/v1"
	"context"
	"roomService/internal/core/domain"
	"roomService/internal/ports"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RoomServer struct {
	roomv1.UnimplementedRoomServiceServer
	svc ports.RoomService
}

func NewRoomServer(svc ports.RoomService) *RoomServer {
	return &RoomServer{svc: svc}
}

func domainStatusToProto(s domain.RequestStatus) roomv1.RequestStatus {
	switch s {
	case domain.RequestStatusPending:
		return roomv1.RequestStatus_REQUEST_STATUS_PENDING
	case domain.RequestStatusAccepted:
		return roomv1.RequestStatus_REQUEST_STATUS_ACCEPTED
	case domain.RequestStatusRejected:
		return roomv1.RequestStatus_REQUEST_STATUS_REJECTED
	default:
		return roomv1.RequestStatus_REQUEST_STATUS_UNSPECIFIED
	}
}

func (s *RoomServer) CreateRoom(ctx context.Context, req *roomv1.CreateRoomRequest) (*roomv1.CreateRoomResponse, error) {
	result, err := s.svc.CreateRoom(ctx, req.GetName(), req.GetIsPrivate(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	return &roomv1.CreateRoomResponse{
		RoomUuid:  result.RoomUUID,
		CreatedAt: timestamppb.New(result.CreatedAt),
	}, nil
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

func (s *RoomServer) CreateJoinRequest(ctx context.Context, req *roomv1.CreateJoinRequestRequest) (*emptypb.Empty, error) {
	if err := s.svc.CreateJoinRequest(ctx, req.GetRoomUuid(), req.GetPublicKey()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *RoomServer) GetJoinRequests(ctx context.Context, req *roomv1.GetJoinRequestsRequest) (*roomv1.GetJoinRequestsResponse, error) {
	requests, err := s.svc.GetJoinRequests(ctx, req.GetRoomUuid())
	if err != nil {
		return nil, err
	}

	protoRequests := make([]*roomv1.JoinRequest, 0, len(requests))
	for _, r := range requests {
		protoRequests = append(protoRequests, &roomv1.JoinRequest{
			RoomUuid:  r.RoomUUID,
			UserUuid:  r.UserUUID,
			Username:  r.Username,
			PublicKey: r.PublicKey,
			Status:    domainStatusToProto(r.Status),
			CreatedAt: timestamppb.New(r.CreatedAt),
		})
	}

	return &roomv1.GetJoinRequestsResponse{Requests: protoRequests}, nil
}

func (s *RoomServer) RespondToJoinRequest(ctx context.Context, req *roomv1.RespondToJoinRequestRequest) (*emptypb.Empty, error) {
	decision := domain.RequestStatus(req.GetDecision())
	if err := s.svc.RespondToJoinRequest(ctx, req.GetRoomUuid(), req.GetUserUuid(), decision); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *RoomServer) MarkAsRead(_ context.Context, _ *roomv1.MarkAsReadRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "MarkAsRead not implemented")
}

func (s *RoomServer) GetUnreadCount(_ context.Context, _ *roomv1.GetUnreadCountRequest) (*roomv1.GetUnreadCountResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetUnreadCount not implemented")
}

func (s *RoomServer) GetAllowedRooms(ctx context.Context, req *roomv1.GetAllowedRoomsRequest) (*roomv1.GetAllowedRoomsResponse, error) {
	roomUUIDs, err := s.svc.GetAllowedRooms(ctx, req.GetUserUuid())
	if err != nil {
		return nil, err
	}
	return &roomv1.GetAllowedRoomsResponse{RoomUuids: roomUUIDs}, nil
}
