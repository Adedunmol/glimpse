package service

// import (
// 	"context"

// 	"github.com/Adedunmol/glimpse/internal/lib/fcm"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type NotificationService struct {
// 	db        *pgxpool.Pool
// 	fcmClient *fcm.FCMClient
// }

// func NewNotificationService(db *pgxpool.Pool, fcm *fcm.FCMClient) *NotificationService {
// 	return &NotificationService{
// 		db:        db,
// 		fcmClient: fcm,
// 	}
// }

// func (n *NotificationService) SendToUser(ctx context.Context, userID, title, message string) error {
// 	devices, err := n.deviceRepo.GetUserTokens(ctx, userID)
// 	if err != nil {
// 		return err
// 	}

// 	var tokens []string
// 	for _, device := range devices {
// 		tokens = append(tokens, device.PushToken)
// 	}

// 	failedTokens, err := n.fcmClient.SendToUser(ctx, tokens, title, message)
// 	if err != nil {
// 		return err
// 	}

// 	if err = n.deviceRepo.DeleteBulkDevice(ctx, userID, failedTokens); err != nil {
// 		return nil
// 	}

// 	return nil
// }
