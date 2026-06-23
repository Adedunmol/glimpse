package notification

import (
	"context"
	"errors"

	"github.com/Adedunmol/glimpse/internal/model/user_device"
	"github.com/rs/zerolog"
)

var EmptyFcmTokenErr error = errors.New("fcm token is empty")

type DeviceService struct {
	deviceRepo *DeviceRepository
}

func NewDeviceService(deviceRepo *DeviceRepository) *DeviceService {
	return &DeviceService{
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
