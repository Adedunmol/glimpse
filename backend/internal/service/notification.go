package service

import (
	"context"

	"github.com/Adedunmol/glimpse/internal/lib/fcm"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
)

type NotificationService struct {
	server     *server.Server
	deviceRepo *repository.DeviceRepository
	fcmClient  *fcm.FCMClient
}

func NewNotificationService(server *server.Server, repo *repository.DeviceRepository, fcm *fcm.FCMClient) *NotificationService {
	return &NotificationService{
		server:     server,
		deviceRepo: repo,
		fcmClient:  fcm,
	}
}

func (n *NotificationService) SendToUser(ctx context.Context, userID, title, message string) error {
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
