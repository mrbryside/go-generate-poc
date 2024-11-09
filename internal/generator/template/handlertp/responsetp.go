package handlertp

const SuccessJsonResponseTemplate = `return ctx.JSON(#httpStatus#, dto.#responseName#)`

const SuccessNoContentResponseTemplate = `return ctx.NoContent(#httpStatus#)`
