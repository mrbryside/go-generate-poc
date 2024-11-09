package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

func (h Handler) CreateProducts(ctx echo.Context) error {
	var req createProductsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, createProductsValidationError(err))
	}

	return ctx.JSON(http.StatusOK, createProductsResponse{})
}
