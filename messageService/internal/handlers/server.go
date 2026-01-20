package handlers

import (
	messagev2 "backend/proto/message/v2"
	"context"
	"messageService/internal/models"
	"messageService/internal/service"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageServer struct {
	service *service.ServiceLayer
	messagev2.UnimplementedMessageServiceServer
}

func NewMessageServer(serviceLayer *service.ServiceLayer) *MessageServer {
	return &MessageServer{
		service: serviceLayer,
	}
}

func (m MessageServer) SendMessage(_ context.Context, req *messagev2.SendMessageRequest) (_ *messagev2.SendMessageResponse, _ error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "wrong user id: %v :c", err)
	}
	created_at := time.Now()
	msg := models.Message{
		Username:  req.Username,
		UserId:    userId,
		Id:        uuid.New(),
		Message:   req.Content,
		Timestamp: created_at.Unix(),
	}

	if err = m.service.SaveMessage(msg); err != nil {
		return nil, err
	}

	return &messagev2.SendMessageResponse{
		MessageId: msg.Id.String(),
		CreatedAt: timestamppb.New(created_at),
	}, nil
}
func (m MessageServer) EditMessage(_ context.Context, req *messagev2.EditMessageRequest) (_ *emptypb.Empty, _ error) {
	msg := models.EditMessage{
		MessageId: req.MessageId,
		Content:   req.Content,
		UserId:    req.UserId,
	}

	if err := m.service.EditMessage(msg); err != nil {
		return nil, err
	}
	return nil, nil
}
func (m MessageServer) DeleteMessage(_ context.Context, req *messagev2.DeleteMessageRequest) (_ *emptypb.Empty, _ error) {
	msg := models.DeleteMessage{
		UserId:    req.UserId,
		MessageId: req.MessageId,
	}
	if err := m.service.DeleteMessage(msg); err != nil {
		return nil, err
	}
	return nil, nil
}
func (m MessageServer) GetHistory(_ context.Context, req *messagev2.GetHistoryRequest) (_ *messagev2.GetHistoryResponse, _ error) {
	normalMessages := m.service.GetMessages(int(req.Limit))
	protoMessages := make([]*messagev2.Message, 0)
	for _, msg := range normalMessages {
		protoMessages = append(protoMessages, msg.ToProto())
	}
	return &messagev2.GetHistoryResponse{
		Messages: protoMessages,
	}, nil
}
