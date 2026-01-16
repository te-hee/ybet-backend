package repository

import (
	v1 "backend/proto/message/v1"
	"context"
	"errors"
	"fmt"
	"gateway/config"
	"gateway/internal/model"
	"log"

	"google.golang.org/grpc/metadata"
)

type RepositoryGrpc struct {
	grpcClient v1.MessageServiceClient
	ctx        context.Context
}

func NewRepositoryGrpc(grpcConn v1.MessageServiceClient, ctx context.Context) *RepositoryGrpc {
	repo := &RepositoryGrpc{
		grpcClient: grpcConn,
		ctx:        ctx,
	}
	return repo
}

func (r *RepositoryGrpc) GetMessageHistory(limit uint32) ([]model.Message, error) {
	request := &v1.GetHistoryRequest{
		Limit: limit,
	}
	md := metadata.Pairs("authorization", "Bearer "+*config.MessageServiceKey)
	ctxWithAuth := metadata.NewOutgoingContext(r.ctx, md)
	response, err := r.grpcClient.GetHistory(ctxWithAuth, request)

	messages := []model.Message{}

	for i := 0; i < len(response.Messages); i++ {
		message, err := model.ProtoToModel(response.Messages[i])
		if err != nil {
			return nil, err
		}
		messages = append(messages, *message)
	}

	return messages, err
}

func (r *RepositoryGrpc) SendMessage(message model.InputMessage) error {
	log.Println(message)
	request := &v1.SendMessageRequest{
		UserId:   message.UserId,
		Username: message.Username,
		Content:  message.Content,
	}
	log.Println(request)

	md := metadata.Pairs("authorization", "Bearer "+*config.MessageServiceKey)
	ctxWithAuth := metadata.NewOutgoingContext(r.ctx, md)

	response, err := r.grpcClient.SendMessage(ctxWithAuth, request)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	if !response.Success {
		return errors.New(response.String())
	}

	return err
}
