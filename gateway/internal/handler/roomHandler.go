package handler

import (
	"gateway/internal/model"
	"gateway/internal/service"
	"gateway/internal/utils"
	"log"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
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

func (h *RoomHandler) HandleRooms(c fiber.Ctx) error {
	switch c.Method() {
	case fiber.MethodGet:
		return h.HandleGetUserRooms(c)
	case fiber.MethodPost:
		return h.HandleCreateRoom(c)
	case fiber.MethodPatch:
		return h.HandleUpdateRoomName(c)
	case fiber.MethodDelete:
		return h.HandleDeleteRoom(c)
	default:
		return utils.WriteJsonErrorWithLog(c, fiber.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *RoomHandler) HandleRoomDetails(c fiber.Ctx) error {
	switch c.Method() {
	case fiber.MethodGet:
		return h.HandleGetRoom(c)
	default:
		return utils.WriteJsonErrorWithLog(c, fiber.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *RoomHandler) HandleMembers(c fiber.Ctx) error{
	switch c.Method() {
	case fiber.MethodGet:
		return h.HandleGetRoomMembers(c)
	default:
		return utils.WriteJsonErrorWithLog(c, fiber.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *RoomHandler) HandleInvites(c fiber.Ctx) error {
	switch c.Method() {
	case fiber.MethodPost:
		return h.HandleCreateInvite(c)
	case fiber.MethodGet:
		return h.HandleGetInvite(c)
	case fiber.MethodDelete:
		return h.HandleDeleteInvite(c)
	default:
		return utils.WriteJsonErrorWithLog(c, fiber.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *RoomHandler) HandleJoinRequests(c fiber.Ctx) error{
	switch c.Method() {
	case fiber.MethodPost:
		return h.HandleCreateJoinRequest(c)
	case fiber.MethodGet:
		return h.HandleGetJoinRequests(c)
	default:
		return utils.WriteJsonErrorWithLog(c, fiber.StatusMethodNotAllowed, "Method not allowed")
	}
}

// ─── Helpers ────────────────────────────────────────────────────────

func handleGRPCError(c fiber.Ctx, err error) error {
	status, errResp := utils.GRPCToHTTPResponse(err)
	log.Println(err)
	return c.Status(status).JSON(errResp)
}

// ─── Room CRUD ──────────────────────────────────────────────────────

func (h *RoomHandler) HandleCreateRoom(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateBody[model.CreateRoomRequest](c)  

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	resp, err := h.service.CreateRoom(token.Raw, *input)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *RoomHandler) HandleGetRoom(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateQuery[model.GetRoomRequest](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}


	resp, err := h.service.GetRoom(token.Raw, input.RoomUUID)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *RoomHandler) HandleUpdateRoomName(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateBody[model.UpdateRoomNameRequest](c) 

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	if err := h.service.UpdateRoomName(token.Raw, *input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *RoomHandler) HandleDeleteRoom(c fiber.Ctx) error {
	token := jwtware.FromContext(c)
	
	input, outputErr := utils.ValidateQuery[model.DeleteRoomRequst](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	if err := h.service.DeleteRoom(token.Raw, input.RoomUUID); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *RoomHandler) HandleGetUserRooms(c fiber.Ctx) error {
	token := jwtware.FromContext(c)


	rooms, err := h.service.GetUserRooms(token.Raw)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(rooms)
}

// ─── Membership ─────────────────────────────────────────────────────

func (h *RoomHandler) HandleGetRoomMembers(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateQuery[model.GetRoomMembersRequest](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	members, err := h.service.GetRoomMembers(token.Raw, input.RoomUUID)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(members)
}

func (h *RoomHandler) HandleLeaveRoom(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateBody[model.LeaveRoomRequest](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	if err := h.service.LeaveRoom(token.Raw, input.RoomUUID); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *RoomHandler) HandleRemoveMember(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateBody[model.RemoveMemberRequest](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	if err := h.service.RemoveMember(token.Raw, *input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// ─── Invite Links ───────────────────────────────────────────────────

func (h *RoomHandler) HandleCreateInvite(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateBody[model.CreateInviteRequest](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	resp, err := h.service.CreateInvite(token.Raw, *input)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *RoomHandler) HandleGetInvite(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateQuery[model.GetInviteRequest](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	resp, err := h.service.GetInvite(token.Raw, input.InvieID)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *RoomHandler) HandleDeleteInvite(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateBody[model.DeleteInviteRequest](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	
	if err := h.service.DeleteInvite(token.Raw, *input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *RoomHandler) HandleJoinViaInvite(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateBody[model.JoinViaInviteRequest](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	if err := h.service.JoinViaInvite(token.Raw, input.InviteID); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// ─── Join Requests ──────────────────────────────────────────────────

func (h *RoomHandler) HandleCreateJoinRequest(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateBody[model.CreateJoinRequestRequest](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	if err := h.service.CreateJoinRequest(token.Raw, *input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *RoomHandler) HandleGetJoinRequests(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateQuery[model.GetJoinReqeustRequest](c)

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	requests, err := h.service.GetJoinRequests(token.Raw, input.RoomUUID)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON((requests))
}

func (h *RoomHandler) HandleRespondToJoinRequest(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateBody[model.RespondToJoinRequestRequest](c) 

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	if err := h.service.RespondToJoinRequest(token.Raw, *input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// ─── Unread Tracking ────────────────────────────────────────────────

func (h *RoomHandler) HandleMarkAsRead(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateBody[model.MarkAsReadRequest](c) 

	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	if err := h.service.MarkAsRead(token.Raw, *input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *RoomHandler) HandleUnreadCount(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	input, outputErr := utils.ValidateQuery[model.UnreadCountReqeust](c)


	if outputErr != nil{
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, outputErr)
	}

	resp, err := h.service.GetUnreadCount(token.Raw, input.RoomUUID)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
