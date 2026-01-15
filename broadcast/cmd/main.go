package main

import (
	"broadcast/config"
	"broadcast/internal/handler"
	messagestream "broadcast/internal/messageStream"
	"broadcast/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/cors"
)

func main() {
	config.InitFlags()

	var nc *nats.Conn
	var err error

	for {
		log.Printf("trying to connect on %v", *config.NatsAddr)
		nc, err = nats.Connect(*config.NatsAddr)
		if err != nil {
			log.Printf("NATS error: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("connected to NATS")
		break
	}
	js, _ := jetstream.New(nc)

	msgChannel := make(chan models.Message, 100)

	messageStream := messagestream.NewMessageStreamClient(js, msgChannel)

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
