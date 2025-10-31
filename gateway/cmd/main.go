package main

import (
	"backend/gateway/internal/handler"
	"backend/gateway/internal/repository"
	"backend/gateway/internal/service"
	"context"
	"flag"
	v1 "github.com/te-hee/ybet-backend/proto/message/v1"
	"log"
	"net/http"

	"github.com/rs/cors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {
	mux := http.NewServeMux()

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

	mux.HandleFunc("/messages", handler.HandleMesseges)
	handlerCORS := cors.Default().Handler(mux)
	if err := http.ListenAndServe(":8080", handlerCORS); err != nil {
		panic(err)
	}
}
