package user_device

import "github.com/go-playground/validator/v10"

type CreateDevicePayload struct {
	DeviceToken string `json:"deviceToken" validate:"required"`
	Platform    string `json:"platform" validate:"required"`
}

func (c *CreateDevicePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
