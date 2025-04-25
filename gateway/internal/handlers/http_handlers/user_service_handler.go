package http_handlers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
	"xlink/common/gen/user_service"
	"xlink/common/logger"
	"xlink/gateway/internal/handlers"
	"xlink/gateway/internal/handlers/helpers"
	"xlink/gateway/internal/schemas"
	"xlink/gateway/internal/services"
)

type UserServiceHandler struct {
	userService *services.UserService
}

func NewUserServiceHandler(userService *services.UserService) *UserServiceHandler {
	return &UserServiceHandler{userService}
}

func (h *UserServiceHandler) CreateUser(ctx *fiber.Ctx) error {
	var body schemas.CreateUserSchema
	if err := ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	falseValue := false

	request := &user_service.CreateUserRequest{
		TgId:    body.TgId,
		IsStaff: &falseValue,
		IsAdmin: &falseValue,
	}

	response, err := h.userService.CreateUser(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response"))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "created user", zap.String("id", response.UserId))
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (h *UserServiceHandler) GetUser(ctx *fiber.Ctx) error {
	userId, err := helpers.ParseUUIDField(ctx, "id")
	if err != nil {
		return helpers.BadRequest(ctx, errors.New("invalid user id: must be a valid uuid"))
	}

	request := &user_service.GetUserRequest{
		UserId: userId.String(),
	}

	response, err := h.userService.GetUser(request)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "received user", zap.String("id", response.UserId))
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *UserServiceHandler) GetUserIDByToken(ctx *fiber.Ctx) error {
	var body schemas.UserIdByTokenSchema
	if err := ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	request := &user_service.GetUserIDByTokenRequest{
		Token: body.Token,
	}

	response, err := h.userService.GetUserIDByToken(request)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "received user ID by token", zap.String("id", response.UserId))
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *UserServiceHandler) GetUserIdByTgId(ctx *fiber.Ctx) error {
	tgIdText := ctx.Params("tg_id")
	if len(tgIdText) == 0 {
		return helpers.BadRequest(ctx, errors.New("invalid tg_id: cannot be empty"))
	}
	tgIdInt, err := strconv.Atoi(tgIdText)
	if err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid tg_id: %v", err))
	}
	tgId := int64(tgIdInt)

	request := &user_service.GetUserIDByTgIDRequest{
		TgId: tgId,
	}

	response, err := h.userService.GetUserIDByTgID(request)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "received user ID by tg_ID", zap.String("id", response.UserId))
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *UserServiceHandler) UpdateUser(ctx *fiber.Ctx) error {
	userId, err := helpers.ParseUUIDField(ctx, "id")
	if err != nil {
		return helpers.BadRequest(ctx, errors.New("invalid user id: must be a valid uuid"))
	}

	var body schemas.UpdateUserSchema
	if err = ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	request := &user_service.UpdateUserRequest{
		UserId:  userId.String(),
		TgId:    &body.TgId,
		IsStaff: nil,
		IsAdmin: nil,
	}

	response, err := h.userService.UpdateUser(request)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(),
			"updated user",
			zap.String("id", userId.String()),
			zap.Bool("status", response.Status),
		)
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (h *UserServiceHandler) CheckToken(ctx *fiber.Ctx) error {
	var body schemas.TokenCheckRequest
	if err := ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	request := &user_service.TokenCheckRequest{
		UserId: body.UserId,
		Token:  body.Token,
	}

	response, err := h.userService.CheckToken(request)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(),
			"checked token for user",
			zap.String("id", body.UserId),
			zap.Bool("status", response.Status),
		)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *UserServiceHandler) RefreshToken(ctx *fiber.Ctx) error {
	var body schemas.RefreshTokenSchema
	if err := ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	request := &user_service.RefreshTokenRequest{
		UserId: body.UserId,
		Token:  body.Token,
	}

	response, err := h.userService.RefreshToken(request)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(),
			"refreshed token for user",
			zap.String("id", body.UserId),
		)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *UserServiceHandler) DeleteUser(ctx *fiber.Ctx) error {
	userId, err := helpers.ParseUUIDField(ctx, "id")
	if err != nil {
		return helpers.BadRequest(ctx, errors.New("invalid user id: must be a valid uuid"))
	}

	request := &user_service.DeleteUserRequest{UserId: userId.String()}

	response, err := h.userService.DeleteUser(request)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	if !response.Status {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "unsuccessful delete operation in user_service")
		return helpers.BadRequest(ctx, errors.New("unsuccessful delete operation in user_service"))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(),
			"deleted user",
			zap.String("id", userId.String()),
		)
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *UserServiceHandler) GetRole(ctx *fiber.Ctx) error {
	userId, err := helpers.ParseUUIDField(ctx, "id")
	if err != nil {
		return helpers.BadRequest(ctx, errors.New("invalid user id: must be a valid uuid"))
	}

	request := &user_service.GetRoleRequest{UserId: userId.String()}

	response, err := h.userService.GetRole(request)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "got roles for user", zap.String("id", userId.String()))
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *UserServiceHandler) CreateUserAdmin(ctx *fiber.Ctx) error {
	var body schemas.CreateUserSchemaAdmin
	if err := ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	request := &user_service.CreateUserRequest{
		TgId:    body.TgId,
		IsStaff: body.IsStaff,
		IsAdmin: body.IsAdmin,
	}

	response, err := h.userService.CreateUser(request)

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response"))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "created user (by admin)", zap.String("id", response.UserId))
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (h *UserServiceHandler) UpdateUserAdmin(ctx *fiber.Ctx) error {
	userId, err := helpers.ParseUUIDField(ctx, "id")
	if err != nil {
		return helpers.BadRequest(ctx, errors.New("invalid user id: must be a valid uuid"))
	}

	var body schemas.UpdateUserSchemaAdmin
	if err = ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	request := &user_service.UpdateUserRequest{
		UserId:  userId.String(),
		TgId:    body.TgId,
		IsStaff: body.IsStaff,
		IsAdmin: body.IsAdmin,
	}

	response, err := h.userService.UpdateUser(request)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(),
			"updated user (by admin)",
			zap.String("id", userId.String()),
			zap.Bool("status", response.Status),
		)
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (h *UserServiceHandler) DeleteUserAdmin(ctx *fiber.Ctx) error {
	userId, err := helpers.ParseUUIDField(ctx, "id")
	if err != nil {
		return helpers.BadRequest(ctx, errors.New("invalid user id: must be a valid uuid"))
	}

	request := &user_service.DeleteUserRequest{
		UserId: userId.String(),
	}

	response, err := h.userService.DeleteUser(request)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(),
			"deleted user (by admin)",
			zap.String("id", userId.String()),
			zap.Bool("status", response.Status),
		)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *UserServiceHandler) Login(ctx *fiber.Ctx) error {
	var body schemas.LoginSchema
	if err := ctx.BodyParser(&body); err != nil {
		return helpers.BadRequest(ctx, fmt.Errorf("invalid body: %v", err))
	}

	responseUserId, err := h.userService.GetUserIDByToken(
		&user_service.GetUserIDByTokenRequest{Token: body.ApiToken},
	)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	responseUser, err := h.userService.GetUser(
		&user_service.GetUserRequest{UserId: responseUserId.UserId},
	)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "logged in user", zap.String("id", responseUserId.UserId))

	return ctx.Status(fiber.StatusOK).JSON(&schemas.LoginResponseSchema{
		UserId:     responseUser.UserId,
		TelegramId: responseUser.TgId,
	})
}

func (h *UserServiceHandler) Profile(ctx *fiber.Ctx) error {
	userIdValue := ctx.Context().Value(handlers.UserIdKey)
	if userIdValue == nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "unauthorized (Use auth middleware before Owner only middleware!!!)"})
	}
	userId := userIdValue.(string)

	response, err := h.userService.GetUser(&user_service.GetUserRequest{UserId: userId})
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
			Error(ctx.UserContext(), "couldn't get user_service response", zap.Error(err))
		return helpers.InternalServerError(ctx,
			fmt.Errorf("couldn't get user_service response: %v", err))
	}

	logger.GetOrCreateLoggerFromCtx(ctx.UserContext()).
		Info(ctx.UserContext(), "got user profile", zap.String("id", userId))

	return ctx.Status(fiber.StatusOK).JSON(response)
}
