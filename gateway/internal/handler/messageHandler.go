package handler

import (
	"encoding/json"
	"gateway/internal/auth"
	"gateway/internal/model"
	"gateway/internal/service"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

type MessageHander struct {
	service *service.MessageService
}

var decoder = schema.NewDecoder() // decoder for url params

var badRequest = model.NewOutputError("Bad Request")

func NewMessageHandler(service *service.MessageService) *MessageHander {
	messageHander := &MessageHander{
		service: service,
	}
	return messageHander
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Error encoding response:", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, model.NewOutputError(message))
}

func (h *MessageHander) HandleMesseges(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	w.Header().Set("Content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.HandleGetMessageHistory(w, r)
	case http.MethodPost:
		h.HandleSendMessage(w, r)
	case http.MethodPatch:
		h.HandleUpdateMessage(w, r)
	case http.MethodDelete:
		h.HandleDeleteMessage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MessageHander) HandleUpdateMessage(w http.ResponseWriter, r *http.Request) {
	var input model.EditMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}
	userId, _ := auth.UserFromContext(r.Context())

	if userId == "" {
		http.Error(w, "Missing user information", http.StatusUnauthorized)
		return
	}

	input.UserId = userId

	err := h.service.EditMessage(input)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutput(true))
}

func (h *MessageHander) HandleDeleteMessage(w http.ResponseWriter, r *http.Request) {
	var input model.DeleteMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad json", http.StatusBadRequest)
		return
	}
	userId, _ := auth.UserFromContext(r.Context())

	if userId == "" {
		http.Error(w, "Missing user information", http.StatusUnauthorized)
		return
	}

	input.UserId = userId

	err := h.service.DeleteMessage(input)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutput(true))

}

func (h *MessageHander) HandleGetMessageHistory(w http.ResponseWriter, r *http.Request) {
	var input model.InputHistory

	if err := decoder.Decode(&input, r.URL.Query()); err != nil {
		writeError(w, http.StatusBadRequest, "Bad request")
		log.Println(err)
		return
	}

	values, _ := url.ParseQuery(r.URL.RawQuery)
	if _, exists := values["limit"]; !exists {
		writeError(w, http.StatusBadRequest, "Bad request")
		log.Println("No `limit` in query")
		return
	}

	if input.Limit < 1 {
		errorMessage := "Invalid `limit` value"
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

	if input.Content == "" {
		errorMessage := "Content was not passed or its empty"
		writeError(w, http.StatusBadRequest, errorMessage)
		log.Println(errorMessage)
		return
	}
	userId, username := auth.UserFromContext(r.Context())
	if userId == "" || username == "" {
		http.Error(w, "user information missing", http.StatusUnauthorized)
		return
	}

	input.UserId = userId
	input.Username = username
	err := h.service.SendMessage(input)

	if err != nil {
		writeError(w, http.StatusBadRequest, "message service error")
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutputSendMessage())
}
