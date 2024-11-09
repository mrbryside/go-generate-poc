package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetProducts(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, getProductsResponse{})
}
