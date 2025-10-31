package main

import (
	v1 "backend/proto/message/v1"
	"broadcast/internal/handler"
	messagestream "broadcast/internal/messageStream"
	"broadcast/internal/models"
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcClient, err := grpc.NewClient(*serverAddr, opts...)
	if err != nil {
		log.Panicf("Failed to create gRPC client: %v", err)
	}
	defer grpcClient.Close()

	client := v1.NewMessageServiceClient(grpcClient)
	log.Println("new message service client")

	msgChannel := make(chan models.Message, 100)

	messageStream := messagestream.NewMessageStreamClient(client, ctx, msgChannel)

	wsHandler := handler.NewWebsocketHandler(msgChannel)

	go func() {
		for {
			err := messageStream.Listen()
			if err != nil {
				log.Printf("Message stream error: %v. Retrying in 5 seconds...", err)
				time.Sleep(5 * time.Second)
			}
		}
	}()

	go wsHandler.BroadcastMessages()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler.WsHandler)
	handlerCORS := cors.Default().Handler(mux)
	log.Println("waiting for conns on :8081")
	if err := http.ListenAndServe(":8081", handlerCORS); err != nil {
		panic(err)
	}
}
