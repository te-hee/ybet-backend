package handler

import (
	"backend/gateway/internal/model"
	"backend/gateway/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	service *service.MessageService
}

var badRequest = model.NewOutputError("Bad Request")

func NewMessageHandler(service *service.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Error encoding response:", err)
	}
}

func writeErr(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, model.NewOutputError(message))
}

func (h *MessageHandler) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		writeErr(w, http.StatusBadRequest, "Missing 'limit' query param")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1{
		writeErr(w, http.StatusBadRequest, "Invalid 'limit' param")
		return
	}

	messages, err := h.service.GetMessageHistory(uint32(limit))

	if err != nil {
		log.Println("Error fetching messages:", err)
		writeErr(w, http.StatusInternalServerError, "Failed reading message history")
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutputGetHistory(messages))
}

func (h *MessageHandler) handlePostMessages(w http.ResponseWriter, r *http.Request) {
	var message model.InputMessage
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		writeErr(w, http.StatusBadRequest, "Invalid JSON body")
		log.Println("Error parsing request:", err)
		return
	}

	if message.Content == "" {
		errorMessage := "Missing or empty 'content' field"
		writeErr(w, http.StatusBadRequest, errorMessage)
		log.Println(errorMessage)
		return
	}

	err := h.service.SendMessage(message.Content)

	if err != nil {
		log.Println("Error sending message:", err)
		writeErr(w, http.StatusInternalServerError, "Failed to send message")
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutputSendMessage())
}

func (h *MessageHandler) HandleMessages(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %v", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	switch r.Method {
	case http.MethodGet:
		h.handleGetMessages(w, r)
	case http.MethodPost:
		h.handlePostMessages(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
