package messagestream

import (
	v1 "backend/proto/message/v1"
	"broadcast/internal/models"
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageStream struct {
	client         v1.MessageServiceClient
	ctx            context.Context
	messageChannel chan models.Message
}

func NewMessageStreamClient(client v1.MessageServiceClient, ctx context.Context, channel chan models.Message) *MessageStream {
	return &MessageStream{
		client:         client,
		ctx:            ctx,
		messageChannel: channel,
	}
}

func (m *MessageStream) Listen() error {
	log.Println("run Listen()")
	empty := &emptypb.Empty{}
	stream, err := m.client.StreamMessages(m.ctx, empty)
	if err != nil {
		log.Printf("error creating stream: %v", err)
		return err
	}
	_, err = m.client.Ready(m.ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	defer func() {
		if err := stream.CloseSend(); err != nil {
			log.Printf("Error closing stream: %v", err)
		}
	}()

	for {
		log.Println("listening for messages")
		receivedMessage, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return err
		}
		msg := models.Message{
			Target: "",
			Type:   "UserMessage",
			Payload: models.UserMessage{
				Id:        receivedMessage.GetUuid(),
				Message:   receivedMessage.GetContent(),
				Timestamp: uint64(receivedMessage.GetTimestamp()),
			},
		}
		log.Printf("got message: %v", msg)

		m.messageChannel <- msg
		log.Println("sent messages thorugh channel")
	}
}
