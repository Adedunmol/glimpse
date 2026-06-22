package service

import (
	"context"
	"errors"

	"github.com/Adedunmol/glimpse/internal/model/user_device"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/rs/zerolog"
)

var EmptyFcmTokenErr error = errors.New("fcm token is empty")

type DeviceService struct {
	server     *server.Server
	deviceRepo *repository.DeviceRepository
}

func NewDeviceService(server *server.Server, deviceRepo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{
		server:     server,
		deviceRepo: deviceRepo,
	}
}

func (d *DeviceService) RegisterDevice(ctx context.Context, logger *zerolog.Logger, userID, fcmToken, platform string) (*user_device.UserDevice, error) {
	if fcmToken == "" {
		logger.Error().Msg("device registration failed: fcm token is empty")
		return nil, EmptyFcmTokenErr
	}

	return d.deviceRepo.UpsertDevice(ctx, userID, fcmToken, platform)
}
