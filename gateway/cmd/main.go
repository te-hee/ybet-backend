package main

import (
	messagev2 "backend/proto/message/v2"
	"context"
	"flag"
	"gateway/config"
	"gateway/internal/auth"
	"gateway/internal/client"
	"gateway/internal/handler"
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

	gprcClient := messagev2.NewMessageServiceClient(grpcClient)

	ctx := context.Background()
	msgClient := client.NewMessageServiceClient(ctx, gprcClient)
	service := service.NewMessageService(msgClient)
	handler := handler.NewMessageHandler(service)

	authClient := auth.NewMinimalClient()
	authService := auth.NewMinimalService(authClient)
	authHandler := auth.NewAuthHandler(authService)

	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/messages", auth.AuthMiddleware(handler.HandleMesseges))
	mux.HandleFunc("/login", authHandler.HandleLogin)
	handlerCORS := cors.Default().Handler(mux)

	log.Println("Running server on port:", *config.GatewayPort)
	log.Println("Connected to message service on:", *config.MessageServiceAddr)

	if err := http.ListenAndServe(":"+*config.GatewayPort, handlerCORS); err != nil {

		panic(err)
	}

}
