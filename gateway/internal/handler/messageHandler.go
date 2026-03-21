package handler

import (
	"gateway/internal/model"
	"gateway/internal/service"
	"gateway/internal/utils"
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
)

type MessageHander struct {
	service *service.MessageService
}



func NewMessageHandler(service *service.MessageService) *MessageHander {
	messageHander := &MessageHander{
		service: service,
	}
	return messageHander
}


func (h *MessageHander) HandleUpdateMessage(c fiber.Ctx) error {
	var input model.EditMessageRequest

	if err := c.Bind().Body(&input); err != nil {
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, "Bad json")
	}

	claims := jwtware.FromContext(c).Claims.(*model.UserClaims)
	userId := claims.Subject

	if userId == "" {
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, "Missing user information")
	}

	input.UserId = userId

	err := h.service.EditMessage(input)
	if err != nil {
		status, errResp := utils.GRPCToHTTPResponse(err)
		return utils.WriteJsonErrorWithLog(c, status, errResp)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *MessageHander) HandleDeleteMessage(c fiber.Ctx) error{
	var input model.DeleteMessageRequest
	if err := c.Bind().Body(&input); err != nil {
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, "Bad json")
	}

	claims := jwtware.FromContext(c).Claims.(*model.UserClaims)
	userId := claims.Subject

	if userId == "" {
		return utils.WriteJsonErrorWithLog(c, fiber.StatusUnauthorized, "Missing user information")
	}

	input.UserId = userId

	err := h.service.DeleteMessage(input)
	if err != nil {
		status, errResp := utils.GRPCToHTTPResponse(err)
		return utils.WriteJsonErrorWithLog(c, status, errResp)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *MessageHander) HandleGetMessageHistory(c fiber.Ctx) error{
	var input model.InputHistory

	if err := c.Bind().Query(&input); err != nil {
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, "Bad request")
	}

	queries := c.Queries()

	if _, exists := queries["limit"]; !exists {
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, "No `limit` in query")
	}

	if input.Limit < 1 {
		return utils.WriteJsonErrorWithLog(c, fiber.StatusBadRequest, "Invalid `limit` value")
	}

	messages, err := h.service.GetMessageHistory(input.Limit)
	if err != nil {
		status, errResp := utils.GRPCToHTTPResponse(err)
		return utils.WriteJsonErrorWithLog(c, status, errResp)
	}

	return c.Status(fiber.StatusOK).JSON(model.NewOutputGetHistory(messages))
}

func (h *MessageHander) HandleSendMessage(c fiber.Ctx) error {
	var input model.InputMessage

	if err := c.Bind().Body(&input); err != nil {
		return utils.WriteErrorMessageWithLog(c, fiber.StatusBadRequest, "Bad request")
	}

	if input.Content == "" {
		return utils.WriteErrorMessageWithLog(c, fiber.StatusBadRequest,  "Content was not passed or is empty")
	}

	user := jwtware.FromContext(c)
	claims := user.Claims.(*model.UserClaims)

	userId := claims.Subject
	username := claims.Username

	if userId == "" || username == "" {
		return utils.WriteErrorMessageWithLog(c, fiber.StatusBadRequest,  "User information missing")
	}

	input.UserId = userId
	input.Username = username
	resp, err := h.service.SendMessage(input)
	if err != nil {
		status, errResp := utils.GRPCToHTTPResponse(err)
		return utils.WriteJsonErrorWithLog(c, status, errResp)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
