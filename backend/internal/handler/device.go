package handler

import (
	"net/http"

	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/Adedunmol/glimpse/internal/model/user_device"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/Adedunmol/glimpse/internal/service"
	"github.com/labstack/echo/v4"
)

type DeviceHandler struct {
	Handler
	deviceService *service.DeviceService
}

func NewDeviceHandler(s *server.Server, deviceService *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		Handler:       NewHandler(s),
		deviceService: deviceService,
	}
}

func (h *DeviceHandler) RegisterDevice(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *user_device.CreateDevicePayload) (*user_device.UserDevice, error) {
			userID := middleware.GetUserID(c)
			logger := middleware.GetLogger(c)

			return h.deviceService.RegisterDevice(c.Request().Context(), logger, userID, payload.DeviceToken, payload.Platform)
		},
		http.StatusCreated,
		&user_device.CreateDevicePayload{},
	)(c)
}
