package http_handlers

import "xlink/gateway/internal/services"

type AnalyticsServiceHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsServiceHandler(analyticsService *services.AnalyticsService) *AnalyticsServiceHandler {
	return &AnalyticsServiceHandler{analyticsService: analyticsService}
}
