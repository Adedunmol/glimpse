package user

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type ClerkEventPayload struct {
	Data       json.RawMessage `json:"data" validate:"required"`
	Object     string          `json:"object" validate:"required"`
	EventType  string          `json:"type" validate:"required"`
	TimeStamp  int64           `json:"timestamp" validate:"required"` // Unix milliseconds
	InstanceID string          `json:"instance_id"`
}

func (p *ClerkEventPayload) Validate() error {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}

type CreateUserDTO struct {
	Email       string
	ClerkUserID string
}

type ClerkUserData struct {
	ID             string              `json:"id"`
	EmailAddresses []ClerkEmailAddress `json:"email_addresses"`
}
type ClerkEmailAddress struct {
	EmailAddress string `json:"email_address"`
}
type ClerkDeletedData struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}
