package handlertp

import (
	"fmt"
)

const (
	ValidationErrorFuncName = "ValidationError"
)

var RequestValidationTemplate = fmt.Sprintf(`
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, %s(err))
	}
`, ValidationErrorFuncName)

var RequestValidationForStatusCodeStyleTemplate = fmt.Sprintf(`validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// you can use this result from the validation error to return the map error message
		_ = %s(err)
		return ctx.JSON(http.StatusBadRequest, %s.%s)
	}
`, ValidationErrorFuncName, DtoFolderAndPackageName, ResponseNameReplaceName)

var ValidationHelperTemplate = fmt.Sprintf(`
package %s`, PackageNameReplaceName) + fmt.Sprintf(`

import (
	"fmt"
	"errors"

	"github.com/go-playground/validator/v10"
)

func %s(err error) map[string]string {`, ValidationErrorFuncName) +
	`var validationErrors validator.ValidationErrors
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
`
