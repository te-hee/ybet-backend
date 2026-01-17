package service

import (
	"encoding/json"
	"fmt"
	"log"
	"messageService/internal/models"
	"messageService/internal/repository"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type ServiceLayer struct {
	repo repository.Repository
	js   jetstream.JetStream
}

func New(repo repository.Repository, js jetstream.JetStream) *ServiceLayer {
	return &ServiceLayer{
		repo: repo,
		js:   js,
	}
}

func (s *ServiceLayer) SaveMessage(message models.Message) error {
	s.repo.SaveMessage(message)

	marshalled, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal json TwT: %v", err)
	}

	t1 := time.Now()
	ackF, err := s.js.PublishAsync("chat.messages.created", marshalled)
	if err != nil {
		return fmt.Errorf("failed to publish on NATS :c : %v", err)
	}

	select {
	case ack := <-ackF.Ok():
		log.Printf("Published msg with sequence number %d on stream %q", ack.Sequence, ack.Stream)
	case err := <-ackF.Err():
		log.Println(err)
	}
	log.Printf("NATS took: %v", time.Since(t1))
	return nil
}

func (s *ServiceLayer) EditMessage(editMessage models.EditMessage) error {
	if err := s.repo.EditMessage(editMessage); err != nil {
		return err
	}
	message := models.NatsEditMessage{
		MessageId: editMessage.MessageId,
		Content:   editMessage.Content,
	}

	marshalled, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal json TwT: %v", err)
	}

	t1 := time.Now()
	ackF, err := s.js.PublishAsync("chat.messages.edited", marshalled)
	if err != nil {
		return fmt.Errorf("failed to publish on NATS :c : %v", err)
	}

	select {
	case ack := <-ackF.Ok():
		log.Printf("Published msg with sequence number %d on stream %q", ack.Sequence, ack.Stream)
	case err := <-ackF.Err():
		return err
	}
	log.Printf("NATS took: %v", time.Since(t1))
	return nil
}

func (s *ServiceLayer) DeleteMessage(deleteMessage models.DeleteMessage) error {
	if err := s.repo.DeleteMessage(deleteMessage); err != nil {
		return err
	}
	message := models.NatsDdeleteMessage{
		MessageId: deleteMessage.MessageId,
	}

	marshalled, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal json TwT: %v", err)
	}

	t1 := time.Now()
	ackF, err := s.js.PublishAsync("chat.messages.deleted", marshalled)
	if err != nil {
		return fmt.Errorf("failed to publish on NATS :c : %v", err)
	}

	select {
	case ack := <-ackF.Ok():
		log.Printf("Published msg with sequence number %d on stream %q", ack.Sequence, ack.Stream)
	case err := <-ackF.Err():
		return err
	}
	log.Printf("NATS took: %v", time.Since(t1))
	return nil
}

func (s *ServiceLayer) GetMessages(limit int) []models.Message {
	return s.repo.GetMessages(limit)
}
