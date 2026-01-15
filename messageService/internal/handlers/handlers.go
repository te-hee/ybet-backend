package handlers

import (
	messagev1 "backend/proto/message/v1"
	"context"
	"encoding/json"
	"log"
	"messageService/config"
	"messageService/internal/models"
	"messageService/internal/service"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageServer struct {
	service *service.ServiceLayer
	messagev1.UnimplementedMessageServiceServer
	messageBroadcast chan models.Message
	js               jetstream.JetStream
}

func NewMessageServer(serviceLayer *service.ServiceLayer, js jetstream.JetStream) *MessageServer {
	return &MessageServer{
		service:          serviceLayer,
		messageBroadcast: make(chan models.Message, *config.CustomBuffer),
		js:               js,
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

	marshalled, err := json.Marshal(msg)
	if err != nil {
		return &messagev1.MessageActionResponse{Success: false}, status.Errorf(codes.Internal, "failed to marshal json: %v", err)
	}

	t1 := time.Now()
	ackF, err := m.js.PublishAsync("chat.messages.created", marshalled)
	if err != nil {
		return &messagev1.MessageActionResponse{Success: false}, status.Errorf(codes.Internal, "failed to publish on NATS jetstream: %v", err)
	}

	select {
	case ack := <-ackF.Ok():
		log.Printf("Published msg with sequence number %d on stream %q", ack.Sequence, ack.Stream)
	case err := <-ackF.Err():
		log.Println(err)
	}
	log.Printf("NATS took: %v", time.Since(t1))

	t2 := time.Now()
	m.messageBroadcast <- msg
	log.Printf("Channel send took: %v", time.Since(t2))
	log.Println(m.messageBroadcast)

	return &messagev1.MessageActionResponse{
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
