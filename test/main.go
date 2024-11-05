package main

import (
	"github.com/labstack/echo/v4"

	"github.com/mrbryside/go-generate/test/handler"
)

func main() {
	e := echo.New()
	h := handler.NewHandler()
	handler.RegisterRoutes(e, h)

	e.Logger.Fatal(e.Start(":8080"))
}
