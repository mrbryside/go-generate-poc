package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

type CreateProductHandler struct {
}

func (h CreateProductHandler) CreateProduct(ctx echo.Context) error {
	var req createProductRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, createProductBadRequestResponse{})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// you can use this result from the validation error to return the map error message
		_ = createProductValidationError(err)
		return ctx.JSON(http.StatusBadRequest, createProductBadRequestResponse{})
	}

	return ctx.JSON(http.StatusOK, createProductOKResponse{})
}
