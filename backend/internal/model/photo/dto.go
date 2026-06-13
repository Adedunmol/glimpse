package photo

import "github.com/go-playground/validator/v10"

type CreatePhotosPayload struct {
	UploadID string `param:"uploadId" validate:"required,uuid"`
	Files    []struct {
		Name string `json:"name" validate:"required,min=1,max=200"`
	} `json:"files" validate:"required,min=1,max=200"`
}

func (p *CreatePhotosPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// -----------------------------------------------------------------------

type DeletePhotosPayload struct {
	UploadID string   `param:"uploadId" validate:"required,uuid"`
	FileID   []string `json:"fileId" validate:"required,uuid"`
}

func (p *DeletePhotosPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// -----------------------------------------------------------------------

type CompletePhotosPayload struct {
	UploadID string `param:"uploadId" validate:"required,uuid"`
	Files    []struct {
		Key string `json:"key" validate:"required"`
	} `json:"files" validate:"required,min=1,max=200"`
}

func (p *CompletePhotosPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
