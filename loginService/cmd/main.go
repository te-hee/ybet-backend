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


func main(){
	err := config.LoadConfig()
	if err != nil{
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", *config.ListenPort)

	if err != nil{
		log.Fatal(err)
	}
	
	grpcServer := grpc.NewServer()

	storage := storage.NewMemoryStorage()
	service := service.NewLogInService(storage, time.Duration(*config.AuthTokenDuration*uint(time.Minute)), time.Duration(*config.RefreshTokenDuration*uint(time.Minute)), *config.Issuer, *config.JwtKey)
	server := server.NewLogInServer(*service, *config.PasswordSalt)

	loginv1.RegisterLoginServiceServer(grpcServer, server)	

	err = grpcServer.Serve(lis)
	if err != nil{
		log.Fatal(err)	
	}
}
