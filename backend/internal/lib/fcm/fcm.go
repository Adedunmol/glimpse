package fcm

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCMClient struct {
	client *messaging.Client
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

func (f *FCMClient) SendToUser(ctx context.Context, deviceTokens []string, title, message string) ([]string, error) {
	var failedTokens []string

	msg := &messaging.MulticastMessage{
		Tokens: deviceTokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
	}

	resp, err := f.client.SendEachForMulticast(ctx, msg)
	if err != nil {
		return nil, err
	}

	for i, r := range resp.Responses {
		if !r.Success && messaging.IsUnregistered(r.Error) {
			failedTokens = append(failedTokens, deviceTokens[i])
		}
	}
	return failedTokens, nil
}
