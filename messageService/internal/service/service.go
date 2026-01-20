package service

import (
	"messageService/internal/dispatcher"
	"messageService/internal/models"
	"messageService/internal/repository"
)

type ServiceLayer struct {
	repo       repository.Repository
	dispatcher *dispatcher.Dispatcher
}

func New(repo repository.Repository, dispatcher *dispatcher.Dispatcher) *ServiceLayer {
	return &ServiceLayer{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

func (s *ServiceLayer) SaveMessage(message models.Message) error {
	s.repo.SaveMessage(message)

	s.dispatcher.Publish("chat.messages.created", message)
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

	s.dispatcher.Publish("chat.messages.edited", message)
	return nil
}

func (s *ServiceLayer) DeleteMessage(deleteMessage models.DeleteMessage) error {
	if err := s.repo.DeleteMessage(deleteMessage); err != nil {
		return err
	}
	message := models.NatsDdeleteMessage{
		MessageId: deleteMessage.MessageId,
	}

	s.dispatcher.Publish("chat.messages.deleted", message)
	return nil
}

func (s *ServiceLayer) GetMessages(limit int) []models.Message {
	return s.repo.GetMessages(limit)
}
