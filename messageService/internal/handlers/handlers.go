package handlers

import (
	messagev1 "backend/proto/message/v1"
	v1 "backend/proto/message/v1"
	"context"
	"log"
	"messageService/config"
	"messageService/internal/models"
	"messageService/internal/service"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageServer struct {
	service *service.ServiceLayer
	messagev1.UnimplementedMessageServiceServer
	messageBroadcast chan models.Message
}

func NewMessageServer(serviceLayer *service.ServiceLayer) *MessageServer {
	return &MessageServer{
		service:          serviceLayer,
		messageBroadcast: make(chan models.Message, *config.CustomBuffer),
	}
}

func (m *MessageServer) SendMessage(ctx context.Context, req *messagev1.SendMessageRequest) (*messagev1.MessageActionResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return &messagev1.MessageActionResponse{Success: false}, status.Errorf(codes.InvalidArgument, "failed parsing user id")
	}
	msg := models.Message{
		Username:  req.Username,
		UserId:    userId,
		Id:        uuid.New(),
		Message:   req.Content,
		Timestamp: time.Now().Unix(),
	}

	m.service.SaveMessage(msg)

	return &messagev1.MessageActionResponse{
		Success: true,
	}, nil
}

func (m *MessageServer) EditMessage(_ context.Context, req *v1.EditMessageRequest) (_ *v1.MessageActionResponse, _ error) {
	msg := models.EditMessage{
		UserId:    req.UserId,
		MessageId: req.MessageId,
		Content:   req.Content,
	}
	if err := m.service.EditMessage(msg); err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.PermissionDenied, "failed to edit the message")
	}
	return &v1.MessageActionResponse{
		Success: true,
	}, nil
}
func (m *MessageServer) DeleteMessage(_ context.Context, req *v1.DeleteMessageRequest) (_ *v1.MessageActionResponse, _ error) {
	msg := models.DeleteMessage{
		UserId:    req.UserId,
		MessageId: req.MessageId,
	}
	if err := m.service.DeleteMessage(msg); err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.PermissionDenied, "failed to delete the message")
	}
	return &v1.MessageActionResponse{
		Success: true,
	}, nil
}

func (m *MessageServer) GetHistory(_ context.Context, req *messagev1.GetHistoryRequest) (*messagev1.GetHistoryResponse, error) {
	limit := int(req.GetLimit())

	messages := m.service.GetMessages(limit)
	resp := make([]*messagev1.Message, 0)
	for _, v := range messages {
		resp = append(resp, v.ToProto())
	}
	log.Println(resp)

	return &messagev1.GetHistoryResponse{
		Messages: resp,
	}, nil
}

func (m *MessageServer) StreamMessages(_ *emptypb.Empty, stream grpc.ServerStreamingServer[messagev1.Message]) error {
	log.Println("new client connected")
	for {
		msg := <-m.messageBroadcast
		if err := stream.Send(msg.ToProto()); err != nil {
			return err
		}
		log.Println("sent message to ws")
	}
}

func (m *MessageServer) Ready(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
