package main

import (
	v1 "backend/proto/message/v1"
	"context"
	"flag"
	"gateway/config"
	"gateway/internal/auth"
	"gateway/internal/handler"
	"gateway/internal/repository"
	"gateway/internal/service"
	"log"
	"net/http"

	"github.com/rs/cors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config.InitConfig()

	mux := http.NewServeMux()

	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	grpcClient, err := grpc.NewClient(*config.MessageServiceAddr, opts...)

	if err != nil {
		log.Panic(err)
	}

	client := v1.NewMessageServiceClient(grpcClient)

	ctx := context.Background()
	repo := repository.NewRepositoryGrpc(client, ctx)
	service := service.NewMessageService(repo)
	handler := handler.NewMessageHandler(service)

	authClient := auth.NewMinimalClient()
	authService := auth.NewMinimalService(authClient)
	authHandler := auth.NewAuthHandler(authService)

	mux.HandleFunc("/messages", auth.AuthMiddleware(handler.HandleMesseges))
	mux.HandleFunc("/login", authHandler.HandleLogin)
	handlerCORS := cors.Default().Handler(mux)

	log.Println("Running server on port:", *config.GatewayPort)
	log.Println("Connected to message service on:", *config.MessageServiceAddr)

	if err := http.ListenAndServe(":"+*config.GatewayPort, handlerCORS); err != nil {

		panic(err)
	}

}
