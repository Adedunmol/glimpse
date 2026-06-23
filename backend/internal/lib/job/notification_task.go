package job

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TaskNotifyLinkCreated  = "notify:link.created"
	TaskNotifyUploadDone   = "notify:upload.completed"
	TaskNotifyClusterReady = "notify:cluster.ready"
)

type NotificationPayload struct {
	UserID  string `json:"userId"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func NewNotificationTask(taskType, userID, title, message string) (*asynq.Task, error) {
	payload, err := json.Marshal(NotificationPayload{
		UserID:  userID,
		Title:   title,
		Message: message,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(taskType, payload,
		asynq.MaxRetry(3),
		asynq.Queue(CriticalPriority),
		asynq.Timeout(30*time.Second)), nil
}
