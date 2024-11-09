package main

import (
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"

	_ "github.com/mrbryside/go-generate/test/docs"
	"github.com/mrbryside/go-generate/test/handler"
)

// @title Service Name
// @version 1.0
// @description Description of the service
// @BasePath  /api/v1
// @tag.name API For Frontend
// @tag.description API For Frontend only. Please aware that this API keeps changing by new requirements.
func main() {
	e := echo.New()
	h := handler.NewHandler()

	api := e.Group("/api/v1")
	handler.RegisterRoutes(api, h)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
