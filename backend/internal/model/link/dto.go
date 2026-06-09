package link

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ----------------------------------------------------------------------------------------------

type CreateLinkPayload struct {
	ClusterID           uuid.UUID  `param:"clusterId" validate:"required,uuid"`
	Token               string     `json:"token" validate:"required,min=1"`
	IsPasswordProtected *bool      `json:"isPasswordProtected"`
	PasswordHash        *string    `json:"passwordHash" validate:"omitempty,min=1"`
	ExpiresAt           *time.Time `json:"expiresAt"`
	IsActive            *bool      `json:"isActive"`
}

func (p *CreateLinkPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ----------------------------------------------------------------------------------------------

type UpdateLinkPayload struct {
	ID                  uuid.UUID  `param:"id" validate:"required,uuid"`
	ClusterID           *uuid.UUID `param:"clusterId" validate:"omitempty,uuid"`
	Token               *string    `json:"token" validate:"required,min=1"`
	IsPasswordProtected *bool      `json:"isPasswordProtected"`
	PasswordHash        *string    `json:"passwordHash" validate:"omitempty,min=1"`
	IsActive            *bool      `json:"isActive"`
	ExpiresAt           *time.Time `json:"expiresAt"`
}

func (p *UpdateLinkPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ----------------------------------------------------------------------------------------------

type GetLinkByClusterIDPayload struct {
	ClusterID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *GetLinkByClusterIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ----------------------------------------------------------------------------------------------

type GetLinksQuery struct {
	Page   *int    `query:"page" validate:"omitempty,min=1"`
	Limit  *int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort   *string `query:"sort" validate:"omitempty,oneof=created_at updated_at name"`
	Order  *string `query:"order" validate:"omitempty,oneof=asc desc"`
	Search *string `query:"search" validate:"omitempty,min=1"`
}

func (q *GetLinksQuery) Validate() error {
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

type DeleteLinkPayload struct {
	ClusterID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeleteLinkPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
