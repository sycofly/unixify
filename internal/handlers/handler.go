package handlers

import (
	"github.com/home/unixify/internal/service"
	"github.com/sirupsen/logrus"
)

// Handler handles HTTP requests
type Handler struct {
	services *service.Services
	logger   *logrus.Logger
}

// NewHandler creates a new handler
func NewHandler(services *service.Services, logger *logrus.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}