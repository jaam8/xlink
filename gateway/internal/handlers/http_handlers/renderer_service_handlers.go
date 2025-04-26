package http_handlers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
	"xlink/gateway/internal/handlers/helpers"
	"xlink/gateway/internal/services"
)

type RendererServiceHandler struct {
	rendererService *services.RendererService
}

func NewRendererServiceHandler(rendererService *services.RendererService) *RendererServiceHandler {
	return &RendererServiceHandler{rendererService: rendererService}
}

func (h *RendererServiceHandler) Image(ctx *fiber.Ctx) error {
	var err error

	var shortLink string
	shortLink, err = helpers.ParseNonEmptyStringField(ctx, "shortLink")
	if err != nil {
		return helpers.BadRequest(ctx, errors.New("'shortLink' is required"))
	}

	var param string
	param, err = helpers.ParseNonEmptyStringField(ctx, "param")
	if err != nil {
		return helpers.BadRequest(ctx, errors.New("'param' is required"))
	}

	var startDate time.Time
	startDate, err = helpers.ParseDateField(ctx, "start_date")
	if err != nil {
		return helpers.InvalidDateBadRequest(ctx, "start_date")
	}

	var endDate time.Time
	endDate, err = helpers.ParseDateField(ctx, "end_date")
	if err != nil {
		return helpers.InvalidDateBadRequest(ctx, "end_date")
	}

	var response []byte
	response, err = h.rendererService.Generate(shortLink, param, startDate, endDate)
	if err != nil {
		return helpers.InternalServerError(ctx, fmt.Errorf("couldn't get renderer response: %w", err))
	}

	return ctx.Status(fiber.StatusOK).Send(response)
}
