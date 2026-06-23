package notification

import (
	"context"
	"fmt"

	"github.com/Adedunmol/glimpse/internal/lib/fcm"
)

type NotificationService struct {
	deviceRepo *DeviceRepository
	fcmClient  *fcm.FCMClient
}

func NewNotificationService(deviceRepo *DeviceRepository, fcm *fcm.FCMClient) *NotificationService {
	return &NotificationService{
		deviceRepo: deviceRepo,
		fcmClient:  fcm,
	}
}

func (n *NotificationService) SendToUser(ctx context.Context, userID, title, message string) error {
	if n.fcmClient == nil {
		return fmt.Errorf("fcm client is uninitialized")
	}
	devices, err := n.deviceRepo.GetUserTokens(ctx, userID)
	if err != nil {
		return err
	}

	var tokens []string
	for _, device := range devices {
		tokens = append(tokens, device.PushToken)
	}

	failedTokens, err := n.fcmClient.SendToUser(ctx, tokens, title, message)
	if err != nil {
		return err
	}

	if err = n.deviceRepo.DeleteBulkDevice(ctx, userID, failedTokens); err != nil {
		return nil
	}

	return nil
}
