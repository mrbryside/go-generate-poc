package handler

import "github.com/labstack/echo/v4"

// this code below is generated by go-generate do not edit manually
// -----------------------------------------------------------------

func RegisterRoutes(e *echo.Echo, h Handler) {
	e.GET("/products", h.GetProducts)
}
