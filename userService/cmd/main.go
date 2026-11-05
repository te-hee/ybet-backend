package main

import (
	userv1 "backend/proto/user/v1"
	"log"
	"userService/config"
	"userService/internal/server"
	"userService/internal/service"
	"userService/internal/storage"
	"net"
	"time"
	"userService/internal/utils"

	"google.golang.org/grpc"
)

func main() {
	config.Load()

	lis, err := net.Listen("tcp", config.Cfg.Server.Port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(utils.AuthInterceptor),
	)

	storage := storage.NewMemoryStorage()
	service := service.NewUserInService(
		storage,
		time.Duration(config.Cfg.Auth.AuthTokenDuration)*time.Minute,
		time.Duration(config.Cfg.Auth.RefreshTokenDuration)*time.Minute,
		config.Cfg.Auth.Issuer,
		config.Cfg.Auth.JwtKey,
	)
	server := server.NewUserInServer(*service, config.Cfg.Auth.PasswordSalt)

	userv1.RegisterUserServiceServer(grpcServer, server)

	log.Printf("user service running on %s", config.Cfg.Server.Port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
