package user_device

import "github.com/Adedunmol/glimpse/internal/model"

type UserDevice struct {
	model.Base

	UserID    string `json:"userId" db:"user_id"`
	PushToken string `json:"pushToken" db:"push_token"`
	Platform  string `json:"platform" db:"platform"`
}
