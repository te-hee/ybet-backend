package main

import (
	messagev2 "backend/proto/message/v2"
	roomv1 "backend/proto/room/v1"
	"context"
	"gateway/config"
	"gateway/internal/auth"
	"gateway/internal/client"
	"gateway/internal/handler"
	"gateway/internal/model"
	"gateway/internal/service"
	"log"
	"time"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config.Load()

	app := fiber.New()
	app.Use(cors.New())	
	api := app.Group("/api")
	
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
	authClient := auth.NewMinimalClient("authClient", time.Minute*10)
	authService := auth.NewMinimalService(authClient)
	authHandler := auth.NewAuthHandler(authService)

	jwtSuccessHandler := func(c fiber.Ctx) error{
		claims := jwtware.FromContext(c).Claims.(*model.UserClaims)
		if claims.TokenType != "Auth"{
			return fiber.ErrUnauthorized
		}
		return c.Next()
	}

	healthHandler := func(c fiber.Ctx) error{
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"status": "ok"})
	}

	// ─── Routes ─────────────────────────────────────────────────────

	v1 := api.Group("/v1")

	v1.Get("/health", healthHandler)
	v1.Post("/login", authHandler.HandleLogin)
	// mux.HandleFunc("/login", authHandler.HandleLogin)

	// ─── Messagges ──────────────────────────────────────────────────

	messages := v1.Group("/messages")


	messages.Use(jwtware.New(
		jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(config.Cfg.Auth.JwtSecret)},
			Extractor: extractors.FromAuthHeader("Bearer"),
			Claims: &model.UserClaims{},
			SuccessHandler: jwtSuccessHandler,
	}))

	// messages.All("/messages", auth.AuthMiddleware(msgHandler.HandleMesseges))
	messages.Get("/", msgHandler.HandleGetMessageHistory)
	messages.Post("/", msgHandler.HandleSendMessage)
	messages.Patch("/", msgHandler.HandleUpdateMessage)
	messages.Delete("/", msgHandler.HandleDeleteMessage)
	// mux.HandleFunc("/messages", auth.AuthMiddleware(msgHandler.HandleMesseges))

	// Room routes

	rooms := v1.Group("/rooms")

	rooms.All("/", auth.AuthMiddleware(roomHandler.HandleRooms))
	rooms.All("/details", auth.AuthMiddleware(roomHandler.HandleRoomDetails))
	rooms.All("/members", auth.AuthMiddleware(roomHandler.HandleMembers))
	rooms.All("/leave", auth.AuthMiddleware(roomHandler.HandleLeaveRoom))
	rooms.All("/remove-member", auth.AuthMiddleware(roomHandler.HandleRemoveMember))
	rooms.All("/invites", auth.AuthMiddleware(roomHandler.HandleInvites))
	rooms.All("/invites/join", auth.AuthMiddleware(roomHandler.HandleJoinViaInvite))
	rooms.All("/join-requests", auth.AuthMiddleware(roomHandler.HandleJoinRequests))
	rooms.All("/join-requests/respond", auth.AuthMiddleware(roomHandler.HandleRespondToJoinRequest))
	rooms.All("/mark-read", auth.AuthMiddleware(roomHandler.HandleMarkAsRead))
	rooms.All("/unread", auth.AuthMiddleware(roomHandler.HandleUnreadCount))

	// handlerCORS := cors.Default().Handler(mux)

	log.Println("Running server on port:", config.Cfg.Server.Port)
	log.Println("Connected to message service on:", config.Cfg.Services.Message.Address)
	log.Println("Connected to room service on:", config.Cfg.Services.Room.Address)

	err = app.Listen(":"+config.Cfg.Server.Port)
	
	if err != nil{
		panic(err)
	}

	// if err := http.ListenAndServe(":"+config.Cfg.Server.Port, handlerCORS); err != nil {
	// 	panic(err)
	// }

}
