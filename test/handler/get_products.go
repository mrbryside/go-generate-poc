package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetProducts(ctx echo.Context) error {
	var req getProductsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, getProductsValidationError(err))
	}

	return ctx.NoContent(http.StatusOK)
}
