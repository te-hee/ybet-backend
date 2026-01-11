package main

import (
	v1 "backend/proto/message/v1"
	"context"
	"flag"
	"gateway/internal/handler"
	"gateway/internal/repository"
	"gateway/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/joho/godotenv"
)


func getEnvs() (string, string){
	err := godotenv.Load()

	var messageServiceIp string;
	var gatewayPort string;

	if err != nil{
		log.Println(".env not found switching do default values")
		messageServiceIp = "localhost:50051"
		gatewayPort = "8080"

	} else{
		log.Println("Running on .env")
		messageServiceIp = os.Getenv("MESSAGE_SERVICE_IP")
		gatewayPort = os.Getenv("GATEWAY_PORT")
	}
	return messageServiceIp, gatewayPort
}

func main() {

	messageServiceIp, gatewayPort := getEnvs()

	var serverAddr = flag.String("addr", messageServiceIp, "The server address in the format of host:port")

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

	log.Println("Running server on port:", gatewayPort)
	log.Println("Connected to message service on:", messageServiceIp)

	if err := http.ListenAndServe(":"+gatewayPort, handlerCORS); err != nil {

		panic(err)
	}


}
