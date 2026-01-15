package messagestream

import (
	"broadcast/internal/models"
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type MessageStream struct {
	ctx            context.Context
	messageChannel chan models.Message
	js             jetstream.JetStream
}

func NewMessageStreamClient(js jetstream.JetStream, channel chan models.Message) *MessageStream {
	return &MessageStream{
		js:             js,
		messageChannel: channel,
	}
}

func (m *MessageStream) Listen() error {
	log.Println("run Listen()")

	_, err := m.js.Conn().QueueSubscribe(
		"chat.messages.created",
		"BROADCAST_QUEUE",
		func(msg *nats.Msg) {
			var recvMsg models.UserMessage
			err := json.Unmarshal(msg.Data, &recvMsg)
			if err != nil {
				log.Printf("error unmarshaling json :c : %v", err)

				msg.Term()
			}
			log.Printf("received a new message: %v", recvMsg)

			m.messageChannel <- models.Message{
				Type:    models.UserMessageType,
				Payload: recvMsg,
			}

			log.Println("message sent to the channel")

			if err := msg.Ack(); err != nil {
				log.Printf("Error sending ACK: %v", err)
			}
		},
	)

	if err != nil {
		log.Printf("error creating subscription: %v", err)
		return err
	}

	log.Println("listening for messages via NATS...")

	select {}
}
