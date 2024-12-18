package handler

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	errs := make(map[string]string)
	if errors.As(err, &validationErrors) {
		for _, validationErr := range validationErrors {
			if validationErr.Tag() == "required" {
				errs[validationErr.Field()] = fmt.Sprintf("%s is required", validationErr.Field())
			}
			if validationErr.Tag() == "email" {
				errs[validationErr.Field()] = fmt.Sprintf("%s must be a valid email", validationErr.Field())
			}
			if validationErr.Tag() == "gte" {
				errs[validationErr.Field()] = fmt.Sprintf("%s must be greater than or equal to %s", validationErr.Field(), validationErr.Param())
			}
			if validationErr.Tag() == "lte" {
				errs[validationErr.Field()] = fmt.Sprintf("%s must be less than or equal to %s", validationErr.Field(), validationErr.Param())
			}
			if validationErr.Tag() == "min" {
				errs[validationErr.Field()] = fmt.Sprintf("%s must be at least %s", validationErr.Field(), validationErr.Param())
			}
			if validationErr.Tag() == "max" {
				errs[validationErr.Field()] = fmt.Sprintf("%s must be at most %s", validationErr.Field(), validationErr.Param())
			}
			if validationErr.Tag() == "len" {
				errs[validationErr.Field()] = fmt.Sprintf("%s must be %s characters long", validationErr.Field(), validationErr.Param())
			}
		}
	}
	return errs
}
