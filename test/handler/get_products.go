package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/mrbryside/go-generate/test/handler/dto"
)

func (h Handler) GetProducts(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, dto.GetProductsResponse{})
}
