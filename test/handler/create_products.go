package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/go-playground/validator/v10"

	"github.com/mrbryside/go-generate/test/handler/dto"
)

func (h Handler) CreateProducts(ctx echo.Context) error {
	var req dto.CreateProductsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.CreateProductsBadRequestResponse{})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// you can use this result from the validation error to return the map error message
		_ = ValidationError(err)
		return ctx.JSON(http.StatusBadRequest, dto.CreateProductsBadRequestResponse{})
	}

	return ctx.JSON(http.StatusOK, dto.CreateProductsOKResponse{})
}
