package handlertp

const TemplateGenerate = `
package handlertp

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// this code below is generated by go-generate do not edit manually
// -----------------------------------------------------------------

`

const Template = `
package handlertp

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

#swaggo#
func (h Handler) #handlerFuncName#(ctx echo.Context) error {
	#requestBind#
	#requestValidation#
	#returnSuccess#
}
`

const HandlerMainTemplate = `
package handlertp

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}
`

//
//func New#handlerFuncName#Router(e *echo.Echo, h #handlerName#) *echo.Echo {
//	e.#method#("#api#", h.#handlerFuncName#)
//	return e
//}
