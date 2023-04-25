package utils

import "github.com/go-playground/validator/v10"

type Error struct {
	FailedField string
	Tag         string
}

var validate = validator.New()

func ValidateStruct(body interface{}) []*Error {
	var errors []*Error
	err := validate.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element Error
			element.FailedField = err.Field()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
	}
	return errors
}
