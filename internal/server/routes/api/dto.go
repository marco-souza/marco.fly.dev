package api

import "github.com/go-playground/validator/v10"

var validate = validator.New(validator.WithRequiredStructEnabled())

type CreateCronInput struct {
	Name    string `json:"name" validate:"required,gte=0,lte=130"`
	Cron    string `json:"cron" validate:"required,gte=9,lte=130"`
	Snippet string `json:"snippet" validate:"required,gte=0"`
}

func (cronInput *CreateCronInput) Validate() error {
	return validate.Struct(cronInput)
}
