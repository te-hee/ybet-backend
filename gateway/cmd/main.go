package main

import (
	messagev2 "backend/proto/message/v2"
	roomv1 "backend/proto/room/v1"
	"context"
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
	config.Load()

	mux := http.NewServeMux()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// ─── Message Service ────────────────────────────────────────────
	grpcClient, err := grpc.NewClient(config.Cfg.Services.Message.Address, opts...)
	if err != nil {
		log.Panic(err)
	}

	gprcClient := messagev2.NewMessageServiceClient(grpcClient)

	ctx := context.Background()
	msgClient := client.NewMessageServiceClient(ctx, gprcClient)
	msgService := service.NewMessageService(msgClient)
	msgHandler := handler.NewMessageHandler(msgService)

	// ─── Room Service ───────────────────────────────────────────────
	roomGrpcConn, err := grpc.NewClient(config.Cfg.Services.Room.Address, opts...)
	if err != nil {
		log.Panic(err)
	}

	roomGrpcClient := roomv1.NewRoomServiceClient(roomGrpcConn)
	roomClient := client.NewRoomServiceClient(ctx, roomGrpcClient)
	roomService := service.NewRoomService(roomClient)
	roomHandler := handler.NewRoomHandler(roomService)

	// ─── Auth ───────────────────────────────────────────────────────
	authClient := auth.NewMinimalClient()
	authService := auth.NewMinimalService(authClient)
	authHandler := auth.NewAuthHandler(authService)

	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}

	// ─── Routes ─────────────────────────────────────────────────────
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/login", authHandler.HandleLogin)
	mux.HandleFunc("/messages", auth.AuthMiddleware(msgHandler.HandleMesseges))

	// Room routes
	mux.HandleFunc("/rooms", auth.AuthMiddleware(roomHandler.HandleRooms))
	mux.HandleFunc("/rooms/details", auth.AuthMiddleware(roomHandler.HandleRoomDetails))
	mux.HandleFunc("/rooms/members", auth.AuthMiddleware(roomHandler.HandleMembers))
	mux.HandleFunc("/rooms/leave", auth.AuthMiddleware(roomHandler.HandleLeaveRoom))
	mux.HandleFunc("/rooms/remove-member", auth.AuthMiddleware(roomHandler.HandleRemoveMember))
	mux.HandleFunc("/rooms/invites", auth.AuthMiddleware(roomHandler.HandleInvites))
	mux.HandleFunc("/rooms/invites/join", auth.AuthMiddleware(roomHandler.HandleJoinViaInvite))
	mux.HandleFunc("/rooms/join-requests", auth.AuthMiddleware(roomHandler.HandleJoinRequests))
	mux.HandleFunc("/rooms/join-requests/respond", auth.AuthMiddleware(roomHandler.HandleRespondToJoinRequest))
	mux.HandleFunc("/rooms/mark-read", auth.AuthMiddleware(roomHandler.HandleMarkAsRead))
	mux.HandleFunc("/rooms/unread", auth.AuthMiddleware(roomHandler.HandleUnreadCount))

	handlerCORS := cors.Default().Handler(mux)

	log.Println("Running server on port:", config.Cfg.Server.Port)
	log.Println("Connected to message service on:", config.Cfg.Services.Message.Address)
	log.Println("Connected to room service on:", config.Cfg.Services.Room.Address)

	if err := http.ListenAndServe(":"+config.Cfg.Server.Port, handlerCORS); err != nil {

		panic(err)
	}

}
