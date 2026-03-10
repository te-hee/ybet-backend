package main

import (
	loginv1 "backend/proto/login/v1"
	"log"
	"loginService/config"
	"loginService/internal/server"
	"loginService/internal/service"
	"loginService/internal/storage"
	"net"
	"time"

	"google.golang.org/grpc"
)

func main() {
	config.Load()

	lis, err := net.Listen("tcp", config.Cfg.Server.Port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	storage := storage.NewMemoryStorage()
	service := service.NewLogInService(
		storage,
		time.Duration(config.Cfg.Auth.AuthTokenDuration)*time.Minute,
		time.Duration(config.Cfg.Auth.RefreshTokenDuration)*time.Minute,
		config.Cfg.Auth.Issuer,
		config.Cfg.Auth.JwtKey,
	)
	server := server.NewLogInServer(*service, config.Cfg.Auth.PasswordSalt)

	loginv1.RegisterLoginServiceServer(grpcServer, server)

	log.Printf("login service running on %s", config.Cfg.Server.Port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
