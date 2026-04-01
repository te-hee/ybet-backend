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
func (h *RoomHandler) MapRooms(r fiber.Router, route string) {
	r.Get(route,h.HandleGetUserRooms)
	r.Post(route, h.HandleCreateRoom)  
	r.Patch(route+":room_uuid<guid>", h.HandleUpdateRoomName)
	r.Delete(route+":room_uuid<guid>", h.HandleDeleteRoom)
}

func (h *RoomHandler) MapRoomDetails(r fiber.Router, route string) {
	r.Get(route,h.HandleGetRoom)
}

func (h *RoomHandler) MapMembers(r fiber.Router, route string) {
	r.Get(route, h.HandleGetRoomMembers)
}

func (h *RoomHandler) MapInvites(r fiber.Router, route string)  {
	r.Get(route+"/:invite_id", h.HandleGetInvite)
	r.Post(route, h.HandleCreateInvite)
	r.Delete(route+"/:room_uuid<guid>/:invite_id", h.HandleDeleteInvite)
}

func (h *RoomHandler) MapJoinRequests(r fiber.Router, router string) {
	r.Get("/", h.HandleGetJoinRequests)
	r.Post("/", h.HandleCreateJoinRequest)
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

	var input model.CreateRoomRequest
	if err := c.Bind().Body(&input); err != nil{
		return err
	}

	resp, err := h.service.CreateRoom(token.Raw, input)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *RoomHandler) HandleGetRoom(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	var input model.GetRoomRequest
	if err := c.Bind().Query(&input); err != nil{
		return err
	}

	resp, err := h.service.GetRoom(token.Raw, input.RoomUUID)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *RoomHandler) HandleUpdateRoomName(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	var input model.UpdateRoomNameRequest
	if err := c.Bind().All(&input); err != nil{
		return err
	}

	if err := h.service.UpdateRoomName(token.Raw, input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *RoomHandler) HandleDeleteRoom(c fiber.Ctx) error {
	token := jwtware.FromContext(c)
	
	var input model.DeleteRoomRequst
	if err := c.Bind().URI(&input); err != nil{
		return err
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

	var input model.GetRoomMembersRequest
	if err := c.Bind().Query(&input); err != nil{
		return err
	}

	members, err := h.service.GetRoomMembers(token.Raw, input.RoomUUID)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(members)
}

func (h *RoomHandler) HandleLeaveRoom(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	var input model.LeaveRoomRequest
	if err := c.Bind().Body(&input); err != nil{
		return err
	}

	if err := h.service.LeaveRoom(token.Raw, input.RoomUUID); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *RoomHandler) HandleRemoveMember(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	var input model.RemoveMemberRequest
	if err := c.Bind().Body(&input); err != nil{
		return err
	}

	if err := h.service.RemoveMember(token.Raw, input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// ─── Invite Links ───────────────────────────────────────────────────

func (h *RoomHandler) HandleCreateInvite(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	var input model.CreateInviteRequest
	if err := c.Bind().Body(&input); err != nil{
		return err
	}

	resp, err := h.service.CreateInvite(token.Raw, input)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *RoomHandler) HandleGetInvite(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	var input model.GetInviteRequest

	if err := c.Bind().URI(&input); err != nil{
		return err
	}
	
	resp, err := h.service.GetInvite(token.Raw, input.InviteID)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *RoomHandler) HandleDeleteInvite(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	var input model.DeleteInviteRequest

	if err := c.Bind().URI(&input); err != nil{
		return err
	}
	

	if err := h.service.DeleteInvite(token.Raw, input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *RoomHandler) HandleJoinViaInvite(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	var input model.JoinViaInviteRequest
	if err := c.Bind().Body(&input); err != nil{
		return err
	}

	if err := h.service.JoinViaInvite(token.Raw, input.InviteID); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// ─── Join Requests ──────────────────────────────────────────────────

func (h *RoomHandler) HandleCreateJoinRequest(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	var input model.CreateJoinRequestRequest
	if err := c.Bind().Body(&input); err != nil{
		return err
	}

	if err := h.service.CreateJoinRequest(token.Raw, input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *RoomHandler) HandleGetJoinRequests(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	var input model.GetJoinReqeustRequest
	if err := c.Bind().Query(&input); err != nil{
		return err
	}

	requests, err := h.service.GetJoinRequests(token.Raw, input.RoomUUID)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON((requests))
}

func (h *RoomHandler) HandleRespondToJoinRequest(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	var input model.RespondToJoinRequestRequest
	if err := c.Bind().Body(&input); err != nil{
		return err
	}

	if err := h.service.RespondToJoinRequest(token.Raw, input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// ─── Unread Tracking ────────────────────────────────────────────────

func (h *RoomHandler) HandleMarkAsRead(c fiber.Ctx) error{
	token := jwtware.FromContext(c)

	var input model.MarkAsReadRequest
	if err := c.Bind().Body(&input); err != nil{
		return err
	}

	if err := h.service.MarkAsRead(token.Raw, input); err != nil {
		return handleGRPCError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *RoomHandler) HandleUnreadCount(c fiber.Ctx) error {
	token := jwtware.FromContext(c)

	var input model.UnreadCountReqeust
	if err := c.Bind().Query(&input); err != nil{
		return err
	}

	resp, err := h.service.GetUnreadCount(token.Raw, input.RoomUUID)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
