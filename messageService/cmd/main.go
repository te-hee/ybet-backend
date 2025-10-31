package main

import (
	messagev1 "backend/proto/message/v1"
	"log"
	"messageService/internal/handlers"
	"messageService/internal/repository"
	"messageService/internal/service"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := loadConfig()
	if err != nil {
		log.Fatalf("error loading env variables ;c: %v", err)
	}
	msgServer := newApp()

	grpcServer := grpc.NewServer()
	messagev1.RegisterMessageServiceServer(grpcServer, msgServer)

	environment := os.Getenv("ENV")
	log.Printf("running on env ^w^: %s", environment)

	if environment == "dev" || environment == "" {
		reflection.Register(grpcServer)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed creating listener TwT: %v", err)
	}

	log.Println("message service running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed serving QwQ: %v", err)
	}
}

func newApp() *handlers.MessageServer {
	repo := repository.NewInMemoryRepo()
	sLayer := service.New(repo)
	server := handlers.NewMessageServer(sLayer)

	return server
}

func loadConfig() error {
	_, err := os.Stat(".env")
	if err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			return err
		}
	}
	return nil
}
