package handlers

import (
	messagev1 "backend/proto/message/v1"
	"context"
	"log"
	"messageService/internal/models"
	"messageService/internal/service"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageServer struct {
	service *service.ServiceLayer
	messagev1.UnimplementedMessageServiceServer
	messageBroadcast chan models.Message
}

func NewMessageServer(serviceLayer *service.ServiceLayer) *MessageServer {
	bufferSize := 100
	if customBuffer := os.Getenv("BUFFER_SIZE"); customBuffer != "" {
		buf, err := strconv.Atoi(customBuffer)
		if err == nil {
			bufferSize = buf
		}
	}
	log.Printf("buffer size :3 : %v", bufferSize)
	return &MessageServer{
		service:          serviceLayer,
		messageBroadcast: make(chan models.Message, bufferSize),
	}
}

func (m *MessageServer) SendMessage(_ context.Context, req *messagev1.SendMessageRequest) (*messagev1.MessageActionResponse, error) {
	msg := models.Message{
		Id:        uuid.New(),
		Message:   req.Content,
		Timestamp: time.Now().Unix(),
	}

	m.service.SaveMessage(msg)

	m.messageBroadcast <- msg
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
