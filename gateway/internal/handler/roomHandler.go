package handler

import (
	"encoding/json"
	"gateway/internal/auth"
	"gateway/internal/model"
	"gateway/internal/service"
	"gateway/internal/utils"
	"log"
	"net/http"
)

type RoomHandler struct {
	service *service.RoomService
}

func NewRoomHandler(service *service.RoomService) *RoomHandler {
	return &RoomHandler{
		service: service,
	}
}

// ─── Route Dispatchers ──────────────────────────────────────────────

func (h *RoomHandler) HandleRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.HandleGetUserRooms(w, r)
	case http.MethodPost:
		h.HandleCreateRoom(w, r)
	case http.MethodPatch:
		h.HandleUpdateRoomName(w, r)
	case http.MethodDelete:
		h.HandleDeleteRoom(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RoomHandler) HandleRoomDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.HandleGetRoom(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RoomHandler) HandleMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.HandleGetRoomMembers(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RoomHandler) HandleInvites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		h.HandleCreateInvite(w, r)
	case http.MethodGet:
		h.HandleGetInvite(w, r)
	case http.MethodDelete:
		h.HandleDeleteInvite(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RoomHandler) HandleJoinRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		h.HandleCreateJoinRequest(w, r)
	case http.MethodGet:
		h.HandleGetJoinRequests(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ─── Helpers ────────────────────────────────────────────────────────

func requireToken(r *http.Request, w http.ResponseWriter) (string, bool) {
	token := auth.RawTokenFromContext(r.Context())
	if token == "" {
		http.Error(w, "Missing user token", http.StatusUnauthorized)
		return "", false
	}
	return token, true
}

func handleGRPCError(w http.ResponseWriter, err error) {
	status, errResp := utils.GRPCToHTTPResponse(err)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResp)
	log.Println(err)
}

// ─── Room CRUD ──────────────────────────────────────────────────────

func (h *RoomHandler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	var input model.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad json")
		return
	}

	if input.Name == "" {
		writeError(w, http.StatusBadRequest, "Room name is required")
		return
	}

	resp, err := h.service.CreateRoom(token, input)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *RoomHandler) HandleGetRoom(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	roomUUID := r.URL.Query().Get("room_uuid")
	if roomUUID == "" {
		writeError(w, http.StatusBadRequest, "room_uuid is required")
		return
	}

	resp, err := h.service.GetRoom(token, roomUUID)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *RoomHandler) HandleUpdateRoomName(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	var input model.UpdateRoomNameRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad json")
		return
	}

	if input.RoomUUID == "" || input.Name == "" {
		writeError(w, http.StatusBadRequest, "room_uuid and name are required")
		return
	}

	if err := h.service.UpdateRoomName(token, input); err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutput(true))
}

func (h *RoomHandler) HandleDeleteRoom(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	roomUUID := r.URL.Query().Get("room_uuid")
	if roomUUID == "" {
		writeError(w, http.StatusBadRequest, "room_uuid is required")
		return
	}

	if err := h.service.DeleteRoom(token, roomUUID); err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutput(true))
}

func (h *RoomHandler) HandleGetUserRooms(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	rooms, err := h.service.GetUserRooms(token)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, rooms)
}

// ─── Membership ─────────────────────────────────────────────────────

func (h *RoomHandler) HandleGetRoomMembers(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	roomUUID := r.URL.Query().Get("room_uuid")
	if roomUUID == "" {
		writeError(w, http.StatusBadRequest, "room_uuid is required")
		return
	}

	members, err := h.service.GetRoomMembers(token, roomUUID)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, members)
}

func (h *RoomHandler) HandleLeaveRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	var input struct {
		RoomUUID string `json:"room_uuid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad json")
		return
	}

	if input.RoomUUID == "" {
		writeError(w, http.StatusBadRequest, "room_uuid is required")
		return
	}

	if err := h.service.LeaveRoom(token, input.RoomUUID); err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutput(true))
}

func (h *RoomHandler) HandleRemoveMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	var input model.RemoveMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad json")
		return
	}

	if input.RoomUUID == "" || input.UserUUID == "" {
		writeError(w, http.StatusBadRequest, "room_uuid and user_uuid are required")
		return
	}

	if err := h.service.RemoveMember(token, input); err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutput(true))
}

// ─── Invite Links ───────────────────────────────────────────────────

func (h *RoomHandler) HandleCreateInvite(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	var input model.CreateInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad json")
		return
	}

	if input.RoomUUID == "" {
		writeError(w, http.StatusBadRequest, "room_uuid is required")
		return
	}

	resp, err := h.service.CreateInvite(token, input)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *RoomHandler) HandleGetInvite(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	inviteID := r.URL.Query().Get("invite_id")
	if inviteID == "" {
		writeError(w, http.StatusBadRequest, "invite_id is required")
		return
	}

	resp, err := h.service.GetInvite(token, inviteID)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *RoomHandler) HandleDeleteInvite(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	var input model.DeleteInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad json")
		return
	}

	if input.InviteID == "" || input.RoomUUID == "" {
		writeError(w, http.StatusBadRequest, "invite_id and room_uuid are required")
		return
	}

	if err := h.service.DeleteInvite(token, input); err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutput(true))
}

func (h *RoomHandler) HandleJoinViaInvite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	var input struct {
		InviteID string `json:"invite_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad json")
		return
	}

	if input.InviteID == "" {
		writeError(w, http.StatusBadRequest, "invite_id is required")
		return
	}

	if err := h.service.JoinViaInvite(token, input.InviteID); err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutput(true))
}

// ─── Join Requests ──────────────────────────────────────────────────

func (h *RoomHandler) HandleCreateJoinRequest(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	var input model.CreateJoinRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad json")
		return
	}

	if input.RoomUUID == "" || input.PublicKey == "" {
		writeError(w, http.StatusBadRequest, "room_uuid and public_key are required")
		return
	}

	if err := h.service.CreateJoinRequest(token, input); err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, model.NewOutput(true))
}

func (h *RoomHandler) HandleGetJoinRequests(w http.ResponseWriter, r *http.Request) {
	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	roomUUID := r.URL.Query().Get("room_uuid")
	if roomUUID == "" {
		writeError(w, http.StatusBadRequest, "room_uuid is required")
		return
	}

	requests, err := h.service.GetJoinRequests(token, roomUUID)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, requests)
}

func (h *RoomHandler) HandleRespondToJoinRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	var input model.RespondToJoinRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad json")
		return
	}

	if input.RoomUUID == "" || input.UserUUID == "" || input.Decision == "" {
		writeError(w, http.StatusBadRequest, "room_uuid, user_uuid, and decision are required")
		return
	}

	if input.Decision != "ACCEPTED" && input.Decision != "REJECTED" {
		writeError(w, http.StatusBadRequest, "decision must be ACCEPTED or REJECTED")
		return
	}

	if err := h.service.RespondToJoinRequest(token, input); err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutput(true))
}

// ─── Unread Tracking ────────────────────────────────────────────────

func (h *RoomHandler) HandleMarkAsRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	var input model.MarkAsReadRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Bad json")
		return
	}

	if input.RoomUUID == "" || input.LastReadMessageID == "" {
		writeError(w, http.StatusBadRequest, "room_uuid and last_read_message_id are required")
		return
	}

	if err := h.service.MarkAsRead(token, input); err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, model.NewOutput(true))
}

func (h *RoomHandler) HandleUnreadCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, ok := requireToken(r, w)
	if !ok {
		return
	}

	roomUUID := r.URL.Query().Get("room_uuid")
	if roomUUID == "" {
		writeError(w, http.StatusBadRequest, "room_uuid is required")
		return
	}

	resp, err := h.service.GetUnreadCount(token, roomUUID)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
