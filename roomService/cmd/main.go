package main

import (
	roomv1 "backend/proto/room/v1"
	"log"
	"net"
	"roomService/config"
	grpcadapter "roomService/internal/adapters/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config.Load()

	roomServer := grpcadapter.NewRoomServer()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcadapter.AuthInterceptor))
	roomv1.RegisterRoomServiceServer(grpcServer, roomServer)

	log.Printf("running on env ^w^: %s", config.Cfg.Env)

	if config.Cfg.Env == "dev" || config.Cfg.Env == "" {
		reflection.Register(grpcServer)
	}

	lis, err := net.Listen("tcp", config.Cfg.Server.Port)
	if err != nil {
		log.Fatalf("failed creating listener TwT: %v", err)
	}

	log.Printf("room service running on %s", config.Cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed serving QwQ: %v", err)
	}
}
