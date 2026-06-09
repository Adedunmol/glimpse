package user

import "github.com/Adedunmol/glimpse/internal/model"

type User struct {
	model.BaseWithCreatedAt
	model.BaseWithUpdatedAt

	ClerkUserID string `json:"clerkUserId" db:"clerk_user_id"`
}
