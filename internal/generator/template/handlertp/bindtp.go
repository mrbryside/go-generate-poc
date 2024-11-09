package handlertp

const BindTemplate = `var req dto.#requestName#
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, #responseName#)
	}
`
