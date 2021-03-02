package validate

import (
	"github.com/go-playground/validator/v10"
)

// V is a shorthand go-playground validation checker using go-playground validator
func V(i interface{}) error {
	if err := validator.New().Struct(i); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			return err
		}
	}
	return nil
}
