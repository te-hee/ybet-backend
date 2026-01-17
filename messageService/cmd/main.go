package main

import (
	messagev1 "backend/proto/message/v1"
	"context"
	"log"
	"messageService/config"
	"messageService/internal/handlers"
	"messageService/internal/repository"
	"messageService/internal/service"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := loadEnvFile()
	log.Println("abababa")
	if err != nil {
		log.Fatalf("error loading env variables ;c: %v", err)
	}
	config.LoadConfig()
	var nc *nats.Conn

	for {
		log.Printf("trying to connect on %v", config.NATSAddress)
		nc, err = nats.Connect(config.NATSAddress)
		if err != nil {
			log.Printf("NATS error: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("connected to NATS")
		break
	}

	msgServer := newApp(nc)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(handlers.AuthInterceptor))
	messagev1.RegisterMessageServiceServer(grpcServer, msgServer)

	log.Printf("running on env ^w^: %s", *config.Env)

	if *config.Env == "dev" || *config.Env == "" {
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

func newApp(nc *nats.Conn) *handlers.MessageServer {
	js, _ := jetstream.New(nc)
	js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name:     "CHAT_MESSAGES",
		Subjects: []string{"chat.messages.>"},
	})

	repo := repository.NewInMemoryRepo()
	sLayer := service.New(repo, js)
	server := handlers.NewMessageServer(sLayer)

	return server
}

func loadEnvFile() error {
	_, err := os.Stat(".env")
	if err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			return err
		}
	}
	return nil
}
