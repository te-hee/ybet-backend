package main

import (
	messagev2 "backend/proto/message/v2"
	roomv1 "backend/proto/room/v1"
	userv1 "backend/proto/user/v1"
	"context"
	"gateway/config"
	"gateway/internal/client"
	"gateway/internal/handler"
	"gateway/internal/model"
	"gateway/internal/service"
	"gateway/internal/utils"
	"log"
	"time"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/gofiber/fiber/v3/middleware/logger"
)



func main() {
	config.Load()

	// ─── Validator ──────────────────────────────────────────────────

	validatorStruct := utils.NewValidateStruct()

	// ─── App Configuration ──────────────────────────────────────────

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.AppErrorHandler,
		StructValidator: validatorStruct,
	})

	app.Use(cors.New())	
	app.Use(logger.New())

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

	// ─── User ───────────────────────────────────────────────────────

	userGrpcConn, err := grpc.NewClient(config.Cfg.Services.User.Address,opts...)
	if err != nil{
		log.Panic(err)
	}

	userCtx := context.Background()
	userClient := client.NewUserClientGrpc(userv1.NewUserServiceClient(userGrpcConn), time.Second*2, userCtx)
	userService := service.NewUserService(userClient)
	userHandler := handler.NewUserHandler(userService)

	jwtSuccessHandler := func(c fiber.Ctx) error{
		claims := jwtware.FromContext(c).Claims.(*model.UserClaims)
		log.Println(claims)
		if claims.TokenType != "Auth" || claims.Username == "" || claims.Subject == ""{
			return fiber.NewError(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
		}
		
		return c.Next()
	}

	healthHandler := func(c fiber.Ctx) error{
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"status": "ok"})
	}

	jwtVerifyMiddleware := jwtware.New(
		jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(config.Cfg.Auth.JwtSecret)},
			Extractor: extractors.FromAuthHeader("Bearer"),
			Claims: &model.UserClaims{},
			SuccessHandler: jwtSuccessHandler,
			ErrorHandler: utils.JwtErrorHandler,
	})

	// ─── Routes ─────────────────────────────────────────────────────

	v1 := api.Group("/v1")

	v1.Get("/health", healthHandler)

	// Login routes

	user := v1.Group("/user")

	user.Post("/sign-in", userHandler.SignIn)
	user.Post("/log-in", userHandler.LogIn)
	user.Post("/auth-token", userHandler.GetNewAuthToken)

	endPointsWithJWTValidation := v1.Group("/")
	endPointsWithJWTValidation.Use(jwtVerifyMiddleware)

	// ─── Messagges ──────────────────────────────────────────────────

	messages := endPointsWithJWTValidation.Group("/messages")

	messages.Get("/", msgHandler.HandleGetMessageHistory)
	messages.Post("/", msgHandler.HandleSendMessage)
	messages.Patch("/:message_id<guid>", msgHandler.HandleUpdateMessage)
	messages.Delete("/:message_id<guid>", msgHandler.HandleDeleteMessage)

	// Room routes

	rooms := endPointsWithJWTValidation.Group("/rooms")
	rooms.Use(jwtVerifyMiddleware)

	roomHandler.MapRooms(rooms, "/")
	roomHandler.MapRoomDetails(rooms, "/details")
	rooms.Post("/leave", roomHandler.HandleLeaveRoom)

	roomHandler.MapMembers(rooms, "/members")
	rooms.Post("/remove-member", roomHandler.HandleRemoveMember)

	roomHandler.MapInvites(rooms, "/invites")
	roomHandler.MapJoinRequests(rooms, "/join-requests")
	rooms.Post("/join-requests/respond", roomHandler.HandleRespondToJoinRequest)
	rooms.Post("/invites/join", roomHandler.HandleJoinViaInvite)

	rooms.Post("/mark-read", roomHandler.HandleMarkAsRead)
	rooms.Get("/unread", roomHandler.HandleUnreadCount)

	log.Println("Running server on port:", config.Cfg.Server.Port)
	log.Println("Connected to message service on:", config.Cfg.Services.Message.Address)
	log.Println("Connected to room service on:", config.Cfg.Services.Room.Address)

	err = app.Listen(":"+config.Cfg.Server.Port)
	
	if err != nil{
		panic(err)
	}

}
