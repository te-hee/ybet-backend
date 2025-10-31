package handler

import (
	"encoding/json"
	"gateway/internal/model"
	"gateway/internal/service"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/schema"
)

type MessageHander struct {
	service *service.MessageService
}

var decoder  = schema.NewDecoder() // decoder for url params

var badRequest = model.NewOutputError("Bad Request")

func NewMessageHandler(service *service.MessageService) *MessageHander {
	messageHander := &MessageHander{
		service: service,
	}
	return messageHander
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Error encoding response:", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, model.NewOutputError(message))
}

func (h *MessageHander) HandleMesseges(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.Body, time.Now())
	w.Header().Set("Content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.HandleGetMessageHistory(w, r)
	case http.MethodPost:
		h.HandleSendMessage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MessageHander) HandleGetMessageHistory(w http.ResponseWriter, r *http.Request) {
	var input model.InputHistory

		if err := decoder.Decode(&input, r.URL.Query()); err != nil {
			writeError(w, http.StatusBadRequest, "Bad request")
			log.Println(err)
			return
		}

	if input.Limit < 1 {
		errorMessage := "Non positiv number passed"
		writeError(w, http.StatusBadRequest, errorMessage)
		log.Println(errorMessage)
		return
	}

	messages, err := h.service.GetMessageHistory(input.Limit)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(badRequest)
		log.Println(err)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutputGetHistory(messages))
}

func (h *MessageHander) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	var input model.InputMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad request")
		log.Println(err)
		return
	}

	if input.Conetnt == "" {
		errorMessage := "Content was not passed or its empty"
		writeError(w, http.StatusBadRequest, errorMessage)
		log.Println(errorMessage)
		return
	}

	err := h.service.SendMessage(input.Conetnt)

	if err != nil {
		writeError(w, http.StatusBadRequest, "Bad request")
		log.Print(err)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutputSendMessage())
}
