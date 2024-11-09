package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

func (h Handler) UpdateProducts(ctx echo.Context) error {
	var req updateProductsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, updateProductsValidationError(err))
	}

	return ctx.JSON(http.StatusOK, updateProductsResponse{})
}
