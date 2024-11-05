package handlertp

const SuccessJsonResponseTemplate = `return ctx.JSON(#httpStatus#, #responseName#)`

const SuccessNoContentResponseTemplate = `return ctx.NoContent(#httpStatus#)`
