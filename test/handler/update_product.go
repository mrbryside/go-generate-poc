package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/go-playground/validator/v10"

	"github.com/mrbryside/go-generate/test/handler/dto"
)

func (h Handler) UpdateProduct(ctx echo.Context) error {
	var req dto.UpdateProductRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ValidationError(err))
	}

	return ctx.JSON(http.StatusOK, dto.UpdateProductResponse{})
}
