package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate() error
}

var validate *validator.Validate

func New(structDateFormat string) error {
	validate = validator.New()

	err := validate.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		if val == "" {
			return true
		}

		_, err := time.Parse(structDateFormat, fl.Field().String())
		return err == nil
	})
	if err != nil {
		return err
	}

	return nil
}
