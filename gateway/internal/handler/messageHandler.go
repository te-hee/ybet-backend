package handler

import (
	"backend/gateway/internal/model"
	"backend/gateway/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type MessageHander struct {
	service *service.MessageService
}

var badRequest = model.OutputError{Output: model.Output{Success: false}, Error: "Bad Request"}

func NewMessageHandler(service *service.MessageService) *MessageHander {
	messageHander := &MessageHander{
		service: service,
	}
	return messageHander
}

func (h *MessageHander) HandleMesseges(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.Body, time.Now())
	switch r.Method {
	case http.MethodGet:

		var input model.InputHistory

		w.Header().Set("Content-type", "application/json")

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(badRequest)
			log.Println(err)
			return
		}

		if input.Limit < 1{
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.OutputError{Output: model.Output{Success: false}, Error: "Non positive number passed"})
			fmt.Println("Non positiv number passed")
			return
		}

		messages, err := h.service.GetMessageHistory(input.Limit)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(badRequest)
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := model.OutputGetHistory{Output: model.Output{Success: true}, Messages: messages}
		json.NewEncoder(w).Encode(response)
		return

	case http.MethodPost:

		badRequest := model.OutputSendMessege{}

		var input model.InputMessage
		w.Header().Set("Content-type", "application/json")
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(badRequest)
			log.Println(err)
			return
		}


		err := h.service.SendMessage(input.Conetnt)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(badRequest)
			log.Print(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := model.OutputSendMessege{Output: model.Output{Success: true}}
		json.NewEncoder(w).Encode(response)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
