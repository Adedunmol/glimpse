package upload

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ----------------------------------------------------------------------------------------------

type CreateUploadPayload struct {
	Name      string     `json:"name" validate:"required,min=1,max=255"`
	Status    *Status    `json:"status" validate:"omitempty,oneof=pending processing done failed"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

func (p *CreateUploadPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ----------------------------------------------------------------------------------------------

type UpdateUploadPayload struct {
	ID        uuid.UUID  `param:"id" validate:"required,uuid"`
	Name      *string    `json:"name" validate:"omitempty,min=1,max=255"`
	Status    *Status    `json:"status" validate:"omitempty,oneof=pending processing done failed"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

func (p *UpdateUploadPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ----------------------------------------------------------------------------------------------

type GetUploadByIDPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *GetUploadByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ----------------------------------------------------------------------------------------------

type GetUploadsQuery struct {
	Page   *int    `query:"page" validate:"omitempty,min=1"`
	Limit  *int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort   *string `query:"sort" validate:"omitempty,oneof=created_at updated_at name"`
	Order  *string `query:"order" validate:"omitempty,oneof=asc desc"`
	Search *string `query:"search" validate:"omitempty,min=1"`
}

func (q *GetUploadsQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(q); err != nil {
		return err
	}

	// set sane defaults
	if q.Page == nil {
		q.Page = new(1)
	}

	if q.Limit == nil {
		q.Limit = new(20)
	}

	if q.Sort == nil {
		q.Sort = new("created_at")
	}

	if q.Order == nil {
		q.Order = new("desc")
	}

	return nil
}

// ----------------------------------------------------------------------------------------------

type DeleteUploadPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeleteUploadPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
