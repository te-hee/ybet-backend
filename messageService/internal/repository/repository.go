package repository

import (
	"messageService/internal/models"
	"sort"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MemoryRepo struct {
	messages map[string]models.Message
	mu       sync.RWMutex
}

func NewInMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		messages: make(map[string]models.Message, 0),
	}
}

func (r *MemoryRepo) SaveMessage(message models.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.messages[message.Id.String()] = message
}
func (r *MemoryRepo) GetMessages(limit int) []models.Message {
	r.mu.Lock()
	defer r.mu.Unlock()

	allMessages := make([]models.Message, 0, len(r.messages))
	for _, v := range r.messages {
		allMessages = append(allMessages, v)
	}

	sort.Slice(allMessages, func(i, j int) bool {
		return allMessages[i].Timestamp < allMessages[j].Timestamp
	})

	if limit > len(allMessages) {
		limit = len(allMessages)
	}

	start := len(allMessages) - limit
	return allMessages[start:]
}

func (r *MemoryRepo) EditMessage(editMessage models.EditMessage) (_ error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var message models.Message
	var ok bool
	if message, ok = r.messages[editMessage.MessageId]; !ok {
		return status.Error(codes.NotFound, "message does not exist :c")
	}

	if message.UserId.String() != editMessage.UserId {
		return status.Error(codes.Unauthenticated, "not your message! >:c")
	}
	message.Message = editMessage.Content
	r.messages[message.Id.String()] = message

	return nil
}
func (r *MemoryRepo) DeleteMessage(deleteMessage models.DeleteMessage) (_ error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var message models.Message
	var ok bool
	if message, ok = r.messages[deleteMessage.MessageId]; !ok {
		return status.Error(codes.NotFound, "message does not exist :c")
	}

	if message.UserId.String() != deleteMessage.UserId {
		return status.Error(codes.Unauthenticated, "not your message! >:c")
	}
	delete(r.messages, message.Id.String())
	return nil
}
