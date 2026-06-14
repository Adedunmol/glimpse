package cluster

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type GetClustersQuery struct {
	Page   *int    `query:"page" validate:"omitempty,min=1"`
	Limit  *int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort   *string `query:"sort" validate:"omitempty,oneof=created_at updated_at label"`
	Order  *string `query:"order" validate:"omitempty,oneof=asc desc"`
	Search *string `query:"search" validate:"omitempty,min=1"`
}

func (q *GetClustersQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(q); err != nil {
		return err
	}

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

type GetClusterByIDPayload struct {
	ID uuid.UUID `param:"clusterId" validate:"required,uuid"`
}

func (p *GetClusterByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
