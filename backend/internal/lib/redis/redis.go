package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const StreamName = "image-tasks"

type Stream struct {
	client *redis.Client
}

func NewStream(client *redis.Client) *Stream {
	return &Stream{client}
}

func (s *Stream) AppendData(ctx context.Context, data map[string]any) error {
	id, err := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: StreamName,
		Values: data,
	}).Result()
	if err != nil {
		return fmt.Errorf("error adding data to stream: %w", err)
	}

	fmt.Printf("data added for id: %s", id)
	return nil
}
