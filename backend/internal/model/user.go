package model

import (
	"github.com/google/uuid"
)

type User struct {
	ClerkIDUser uuid.UUID `json:"userId" db:"user_id"`
	Email       string    `json:"email" db:"email"`
	FCM_Token   string    `json:"-" db:"fcm_token"`
	BaseWithTimeStamp
}
