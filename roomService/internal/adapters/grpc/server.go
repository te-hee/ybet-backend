package grpcadapter

import (
	roomv1 "backend/proto/room/v1"
	"context"
	"roomService/internal/core/domain"
	"roomService/internal/ports"
	"roomService/internal/utils"

	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RoomServer struct {
	roomv1.UnimplementedRoomServiceServer
	service  ports.RoomService
	validate *validator.Validate
}

func NewRoomServer(service ports.RoomService) *RoomServer {
	return &RoomServer{
		service:  service,
		validate: validator.New(),
	}
}

func (s *RoomServer) GetPendingKeys(ctx context.Context, _ *emptypb.Empty) (*roomv1.PendingKeysResponse, error) {
	keys, err := s.service.GetPendingKeys(ctx)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	resp := &roomv1.PendingKeysResponse{}
	for _, k := range keys {
		resp.PendingKeys = append(resp.PendingKeys, &roomv1.PendingKey{
			RoomUuid: k.RoomUUID,
		})
	}
	return resp, nil
}

func (s *RoomServer) AcknowledgeKeyDelivery(ctx context.Context, req *roomv1.AcknowledgeKeyDeliveryRequest) (*emptypb.Empty, error) {
	dto := domain.AcknowledgeKeyDeliveryDTO{
		RoomUUID: req.RoomUuid,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	err := s.service.AcknowledgeKeyDelivery(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *RoomServer) CreateRoom(ctx context.Context, req *roomv1.CreateRoomRequest) (*roomv1.CreateRoomResponse, error) {
	dto := domain.CreateRoomDTO{
		Name:      req.Name,
		IsPrivate: req.IsPrivate,
		GroupID:   req.GroupId,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	res, err := s.service.CreateRoom(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &roomv1.CreateRoomResponse{
		RoomUuid:  res.RoomUUID,
		CreatedAt: timestamppb.New(res.CreatedAt),
	}, nil
}

func (s *RoomServer) GetRoom(ctx context.Context, req *roomv1.GetRoomRequest) (*roomv1.GetRoomResponse, error) {
	dto := domain.GetRoomDTO{
		RoomUUID: req.RoomUuid,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	room, err := s.service.GetRoom(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &roomv1.GetRoomResponse{
		Room: &roomv1.Room{
			RoomUuid:    room.RoomUUID,
			Name:        room.Name,
			AdminId:     room.AdminID,
			IsPrivate:   room.IsPrivate,
			GroupId:     room.GroupID,
			MemberCount: int32(room.MemberCount),
			CreatedAt:   timestamppb.New(room.CreatedAt),
			UpdatedAt:   timestamppb.New(room.UpdatedAt),
		},
	}, nil
}

func (s *RoomServer) UpdateRoomName(ctx context.Context, req *roomv1.UpdateRoomNameRequest) (*emptypb.Empty, error) {
	dto := domain.UpdateRoomNameDTO{
		RoomUUID: req.RoomUuid,
		Name:     req.Name,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	err := s.service.UpdateRoomName(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *RoomServer) DeleteRoom(ctx context.Context, req *roomv1.DeleteRoomRequest) (*emptypb.Empty, error) {
	dto := domain.DeleteRoomDTO{
		RoomUUID: req.RoomUuid,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	err := s.service.DeleteRoom(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *RoomServer) GetUserRooms(ctx context.Context, _ *roomv1.GetUserRoomsRequest) (*roomv1.GetUserRoomsResponse, error) {
	rooms, err := s.service.GetUserRooms(ctx)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	resp := &roomv1.GetUserRoomsResponse{}
	for _, r := range rooms {
		resp.Rooms = append(resp.Rooms, &roomv1.UserRoom{
			RoomUuid:    r.RoomUUID,
			Name:        r.Name,
			IsPrivate:   r.IsPrivate,
			UnreadCount: int32(r.UnreadCount),
			JoinedAt:    timestamppb.New(r.JoinedAt),
			UpdatedAt:   timestamppb.New(r.UpdatedAt),
		})
	}

	return resp, nil
}

func (s *RoomServer) GetRoomMembers(ctx context.Context, req *roomv1.GetRoomMembersRequest) (*roomv1.GetRoomMembersResponse, error) {
	dto := domain.GetRoomMembersDTO{
		RoomUUID: req.RoomUuid,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	members, err := s.service.GetRoomMembers(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	resp := &roomv1.GetRoomMembersResponse{}
	for _, m := range members {
		resp.Members = append(resp.Members, &roomv1.RoomMember{
			UserUuid: m.UserUUID,
			JoinedAt: timestamppb.New(m.JoinedAt),
		})
	}

	return resp, nil
}

func (s *RoomServer) LeaveRoom(ctx context.Context, req *roomv1.LeaveRoomRequest) (*emptypb.Empty, error) {
	dto := domain.LeaveRoomDTO{
		RoomUUID: req.RoomUuid,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	err := s.service.LeaveRoom(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *RoomServer) RemoveMember(ctx context.Context, req *roomv1.RemoveMemberRequest) (*emptypb.Empty, error) {
	dto := domain.RemoveMemberDTO{
		RoomUUID: req.RoomUuid,
		UserUUID: req.UserUuid,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	err := s.service.RemoveMember(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *RoomServer) CreateInvite(ctx context.Context, req *roomv1.CreateInviteRequest) (*roomv1.CreateInviteResponse, error) {
	dto := domain.CreateInviteDTO{
		RoomUUID:  req.RoomUuid,
		UsesLeft:  req.UsesLeft,
		ExpiresAt: req.ExpiresAt,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	res, err := s.service.CreateInvite(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &roomv1.CreateInviteResponse{
		InviteId: res.InviteID,
	}, nil
}

func (s *RoomServer) GetInvite(ctx context.Context, req *roomv1.GetInviteRequest) (*roomv1.GetInviteResponse, error) {
	dto := domain.GetInviteDTO{
		InviteID: req.InviteId,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	invite, err := s.service.GetInvite(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &roomv1.GetInviteResponse{
		Invite: &roomv1.RoomInvite{
			InviteId:  invite.InviteID,
			RoomUuid:  invite.RoomUUID,
			UsesLeft:  int32(invite.UsesLeft),
			ExpiresAt: timestamppb.New(invite.ExpiresAt),
			CreatedAt: timestamppb.New(invite.CreatedAt),
		},
	}, nil
}

func (s *RoomServer) DeleteInvite(ctx context.Context, req *roomv1.DeleteInviteRequest) (*emptypb.Empty, error) {
	dto := domain.DeleteInviteDTO{
		InviteID: req.InviteId,
		RoomUUID: req.RoomUuid,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	err := s.service.DeleteInvite(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *RoomServer) JoinViaInvite(ctx context.Context, req *roomv1.JoinViaInviteRequest) (*emptypb.Empty, error) {
	dto := domain.JoinViaInviteDTO{
		InviteID: req.InviteId,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	err := s.service.JoinViaInvite(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *RoomServer) CreateJoinRequest(ctx context.Context, req *roomv1.CreateJoinRequestRequest) (*emptypb.Empty, error) {
	dto := domain.CreateJoinRequestDTO{
		RoomUUID:  req.RoomUuid,
		PublicKey: req.PublicKey,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	err := s.service.CreateJoinRequest(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *RoomServer) GetJoinRequests(ctx context.Context, req *roomv1.GetJoinRequestsRequest) (*roomv1.GetJoinRequestsResponse, error) {
	dto := domain.GetJoinRequestsDTO{
		RoomUUID: req.RoomUuid,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	reqs, err := s.service.GetJoinRequests(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	resp := &roomv1.GetJoinRequestsResponse{}
	for _, r := range reqs {
		resp.Requests = append(resp.Requests, &roomv1.JoinRequest{
			RoomUuid:  r.RoomUUID,
			UserUuid:  r.UserUUID,
			Username:  r.Username,
			PublicKey: r.PublicKey,
			Status:    roomv1.RequestStatus(r.Status),
			CreatedAt: timestamppb.New(r.CreatedAt),
		})
	}

	return resp, nil
}

func (s *RoomServer) RespondToJoinRequest(ctx context.Context, req *roomv1.RespondToJoinRequestRequest) (*emptypb.Empty, error) {
	dto := domain.RespondToJoinRequestDTO{
		RoomUUID:     req.RoomUuid,
		UserUUID:     req.UserUuid,
		Decision:     domain.RequestStatus(req.Decision),
		EncryptedKey: req.EncryptedKey,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	err := s.service.RespondToJoinRequest(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *RoomServer) MarkAsRead(ctx context.Context, req *roomv1.MarkAsReadRequest) (*emptypb.Empty, error) {
	dto := domain.MarkAsReadDTO{
		RoomUUID:          req.RoomUuid,
		LastReadMessageID: req.LastReadMessageId,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	err := s.service.MarkAsRead(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *RoomServer) GetUnreadCount(ctx context.Context, req *roomv1.GetUnreadCountRequest) (*roomv1.GetUnreadCountResponse, error) {
	dto := domain.GetUnreadCountDTO{
		RoomUUID: req.RoomUuid,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	res, err := s.service.GetUnreadCount(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &roomv1.GetUnreadCountResponse{
		UnreadCount: int32(res.UnreadCount),
	}, nil
}

func (s *RoomServer) GetAllowedRooms(ctx context.Context, req *roomv1.GetAllowedRoomsRequest) (*roomv1.GetAllowedRoomsResponse, error) {
	dto := domain.GetAllowedRoomsDTO{
		UserUUID: req.UserUuid,
	}
	if err := s.validate.Struct(dto); err != nil {
		return nil, utils.AppErrorToGrpcError(domain.NewError(domain.CodeInvalidArgument, "%s", err.Error()))
	}

	rooms, err := s.service.GetAllowedRooms(ctx, dto)
	if err != nil {
		return nil, utils.AppErrorToGrpcError(err)
	}

	return &roomv1.GetAllowedRoomsResponse{
		RoomUuids: rooms,
	}, nil
}
