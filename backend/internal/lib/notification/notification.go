package notification

import (
	"context"
	"fmt"

	"github.com/Adedunmol/glimpse/internal/lib/fcm"
	"github.com/rs/zerolog"
)

type NotificationService struct {
	logger     *zerolog.Logger
	deviceRepo *DeviceRepository
	fcmClient  *fcm.FCMClient
}

func NewNotificationService(logger *zerolog.Logger, deviceRepo *DeviceRepository, fcm *fcm.FCMClient) *NotificationService {
	return &NotificationService{
		logger:     logger,
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

	if len(devices) == 0 {
		n.logger.Info().
			Str("user_id", userID).
			Msg("no registered devices for user, skipping notification")
		return nil
	}

	tokens := make([]string, 0, len(devices))
	for _, device := range devices {
		tokens = append(tokens, device.PushToken)
	}

	result, err := n.fcmClient.SendToUser(ctx, tokens, title, message)
	if err != nil {
		return err
	}

	n.logger.Info().
		Str("user_id", userID).
		Int("devices", len(tokens)).
		Int("delivered", result.SuccessCount).
		Int("failed", result.FailureCount).
		Msg("fcm notification dispatched")

	// FCM reports delivery problems per token rather than as a single error, so
	// these have to be logged explicitly or a send that reached nobody looks
	// identical to one that reached everybody
	for _, failure := range result.Failures {
		n.logger.Error().
			Err(failure.Err).
			Str("user_id", userID).
			Str("device_token", maskToken(failure.Token)).
			Msg("fcm delivery failed for device")
	}

	if len(result.UnregisteredTokens) > 0 {
		n.logger.Info().
			Str("user_id", userID).
			Int("count", len(result.UnregisteredTokens)).
			Msg("pruning unregistered device tokens")

		// pruning is best effort: the notification has already been dispatched,
		// so a cleanup failure must not fail the send
		if err := n.deviceRepo.DeleteBulkDevice(ctx, userID, result.UnregisteredTokens); err != nil {
			n.logger.Error().
				Err(err).
				Str("user_id", userID).
				Msg("failed to prune unregistered device tokens")
		}
	}

	return nil
}

// maskToken keeps enough of a push token to correlate log lines without
// writing the whole credential into the logs.
func maskToken(token string) string {
	if len(token) <= 12 {
		return "***"
	}
	return token[:6] + "..." + token[len(token)-4:]
}
