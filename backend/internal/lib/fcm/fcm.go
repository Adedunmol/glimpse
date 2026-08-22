package fcm

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCMClient struct {
	client *messaging.Client
}

// TokenError pairs a device token with the error FCM returned for it.
type TokenError struct {
	Token string
	Err   error
}

// SendResult reports the per-token outcome of a multicast send. Callers are
// expected to log Failures: FCM reports delivery problems per token rather
// than as a single error, so a send can "succeed" while reaching nobody.
type SendResult struct {
	SuccessCount int
	FailureCount int

	// UnregisteredTokens are tokens FCM reports as no longer registered - the
	// app was uninstalled or the token was rotated. These are safe to delete.
	UnregisteredTokens []string

	// Failures holds every per-token failure, including ones that must not be
	// acted on automatically (quota, sender mismatch, malformed token).
	Failures []TokenError
}

func NewFCMClient(credPath, projectID string) (*FCMClient, error) {
	ctx := context.Background()

	opt := option.WithAuthCredentialsFile(option.ServiceAccount, credPath)
	config := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &FCMClient{
		client: client,
	}, nil
}

func (f *FCMClient) SendToUser(ctx context.Context, deviceTokens []string, title, message string) (*SendResult, error) {
	if len(deviceTokens) == 0 {
		return &SendResult{}, nil
	}

	msg := &messaging.MulticastMessage{
		Tokens: deviceTokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
	}

	resp, err := f.client.SendEachForMulticast(ctx, msg)
	if err != nil {
		// a transport/credential level failure - nothing was delivered
		return nil, fmt.Errorf("fcm multicast send failed: %w", err)
	}

	result := &SendResult{
		SuccessCount: resp.SuccessCount,
		FailureCount: resp.FailureCount,
	}

	for i, r := range resp.Responses {
		if r.Success {
			continue
		}

		result.Failures = append(result.Failures, TokenError{Token: deviceTokens[i], Err: r.Error})

		if messaging.IsUnregistered(r.Error) {
			result.UnregisteredTokens = append(result.UnregisteredTokens, deviceTokens[i])
		}
	}

	return result, nil
}
