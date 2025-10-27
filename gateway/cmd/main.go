package main

import (
	"backend/gateway/internal/handler"
	"backend/gateway/internal/repository"
	"backend/gateway/internal/service"
	v1 "backend/proto/message/v1"
	"context"
	"flag"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {

	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	grpcClient, err := grpc.NewClient(*serverAddr, opts...)

	if err != nil {
		log.Panic(err)
	}

	client := v1.NewMessageServiceClient(grpcClient)

	ctx := context.Background()
	repo := repository.NewRepositoryGrpc(client, ctx)
	service := service.NewMessageService(repo)
	handler := handler.NewMessageHandler(service)

	http.HandleFunc("/messages", handler.HandleMesseges)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
