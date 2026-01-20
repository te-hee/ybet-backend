package client

import (
	messagev2 "backend/proto/message/v2"
	"context"
	"gateway/config"
	"gateway/internal/model"

	"google.golang.org/grpc/metadata"
)

type MessageServiceClient struct {
	ctx        context.Context
	grpcClient messagev2.MessageServiceClient
}

func NewMessageServiceClient(ctx context.Context, grpcClient messagev2.MessageServiceClient) *MessageServiceClient {
	return &MessageServiceClient{
		ctx:        ctx,
		grpcClient: grpcClient,
	}
}

func (c MessageServiceClient) GetMessageHistory(limit uint32) (_ []model.Message, _ error) {
	request := &messagev2.GetHistoryRequest{
		Limit: limit,
	}

	ctxWithAuth := setAuth(c.ctx)

	response, err := c.grpcClient.GetHistory(ctxWithAuth, request)
	messages := make([]model.Message, 0)
	if err != nil {
		return messages, err
	}

	for _, v := range response.Messages {
		messages = append(messages, ProtoToMessage(v))
	}
	return messages, nil
}
func (c MessageServiceClient) SendMessage(message model.InputMessage) (*model.OutputSendMessege, error) {
	request := &messagev2.SendMessageRequest{
		UserId:   message.UserId,
		Username: message.Username,
		Content:  message.Content,
	}
	ctxWithAuth := setAuth(c.ctx)
	response, err := c.grpcClient.SendMessage(ctxWithAuth, request)
	if err != nil {
		return nil, err
	}

	return &model.OutputSendMessege{
		MessageId: response.MessageId,
		Timestamp: response.CreatedAt.AsTime().Unix(),
	}, nil
}
func (c MessageServiceClient) EditMessage(editRequest model.EditMessageRequest) (_ error) {
	request := &messagev2.EditMessageRequest{
		UserId:    editRequest.UserId,
		MessageId: editRequest.MessageId,
		Content:   editRequest.Content,
	}
	ctxWithAuth := setAuth(c.ctx)

	_, err := c.grpcClient.EditMessage(ctxWithAuth, request)
	if err != nil {
		return err
	}
	return nil
}
func (c MessageServiceClient) DeleteMessage(deleteRequest model.DeleteMessageRequest) (_ error) {
	request := &messagev2.DeleteMessageRequest{
		UserId:    deleteRequest.UserId,
		MessageId: deleteRequest.MessageId,
	}
	ctxWithAuth := setAuth(c.ctx)

	_, err := c.grpcClient.DeleteMessage(ctxWithAuth, request)
	if err != nil {
		return err
	}
	return nil
}

func setAuth(ctx context.Context) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+*config.MessageServiceKey)
	ctxWithAuth := metadata.NewOutgoingContext(ctx, md)
	return ctxWithAuth
}

func ProtoToMessage(message *messagev2.Message) model.Message {
	return model.Message{
		MessageId: message.MessageId,
		UserId:    message.UserId,
		Username:  message.Username,
		Timestamp: message.CreatedAt.AsTime().Unix(),
		Content:   message.Content,
	}
}
