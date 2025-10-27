package repository

import (
	"backend/gateway/internal/model"
	"backend/proto/message/v1"
	"context"
	"errors"
)

type RepositoryGrpc struct{
	grpcClient v1.MessageServiceClient
	ctx context.Context
}

func NewRepositoryGrpc(grpcConn v1.MessageServiceClient, ctx context.Context) *RepositoryGrpc{
	repo := &RepositoryGrpc{
		grpcClient: grpcConn,
		ctx: ctx,
	}
	return repo
}

func (r *RepositoryGrpc) GetMessageHistory(limit uint32) ([]model.Message, error){
	request := &v1.GetHistoryRequest{
		Limit: limit,
	}
	response, err := r.grpcClient.GetHistory(r.ctx,request)

	messages := []model.Message{}

	for i:=0; i<len(response.Messages); i++{
 		message, err := model.ProtoToModel(*response.Messages[i])
		if err != nil{
			return nil, err
		}
		messages = append(messages, *message)
	}

	return messages, err
}

func (r* RepositoryGrpc) SendMessage(content string) (error){
	request := &v1.SendMessageRequest{
		Content: content,
	}

	response, err := r.grpcClient.SendMessage(r.ctx, request)

	if !response.Success{
		return errors.New(response.String())
	}

	return err
}