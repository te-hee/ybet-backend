package client

import (
	userv1 "backend/proto/user/v1"
	"context"
	"gateway/config"
	"gateway/internal/model"
	"gateway/internal/utils"
	"time"
)

type UserClientGrpc struct{
	grpcClient userv1.UserServiceClient
	contextWaitTime time.Duration
	ctx context.Context
}



func NewUserClientGrpc(grpcClient userv1.UserServiceClient, contextWaitTime time.Duration, ctx context.Context) *UserClientGrpc{
	return &UserClientGrpc{grpcClient: grpcClient, contextWaitTime: contextWaitTime, ctx: ctx}
}


func (c *UserClientGrpc)SignIn(password string, username string) (*model.SignInResponseV2, error){
	ctx, cancel := context.WithTimeout(c.ctx, c.contextWaitTime)
	defer cancel()

	ctx = utils.SetAuth(ctx, config.Cfg.Services.User.ApiKey)
	
	r, err := c.grpcClient.SignIn(ctx, &userv1.SignInRequest{Password: password, Username: username})	

	if err != nil{
		return nil,  err
	}

	id, err := utils.ParseUUID(r.UserId) 

	if err != nil{
		return nil, err
	}

	return &model.SignInResponseV2{Id: *id, AuthToken: r.AuthToken, RefreshToken: r.RefreshToken}, nil
}

func 	(c *UserClientGrpc)LogIn(password string, username string) (*model.LogInResponseV2, error){
	ctx, cancel := context.WithTimeout(c.ctx, c.contextWaitTime)
	defer cancel()

	ctx = utils.SetAuth(ctx, config.Cfg.Services.User.ApiKey)

	r, err := c.grpcClient.LogIn(ctx, &userv1.LogInRequest{Password: password, Username: username})
	
	if err != nil{
		return nil, err
	}

	id, err := utils.ParseUUID(r.UserId)

	if err != nil {
		return nil, err
	}

	return &model.LogInResponseV2{Id: *id, AuthToken: r.AuthToken, RefreshToken: r.RefreshToken}, nil
}

func (c *UserClientGrpc)GetNewAuthToken(refreshToken string) (*string , error){
	ctx, cancel := context.WithTimeout(c.ctx, c.contextWaitTime)
	defer cancel()

	ctx = utils.SetAuth(ctx, config.Cfg.Services.User.ApiKey)

	r, err := c.grpcClient.GetNewAuthToken(ctx, &userv1.GetNewAuthTokenRequest{RefreshToken: refreshToken})

	if err != nil{
		return nil, err
	}

	return &r.AuthToken, nil
}
