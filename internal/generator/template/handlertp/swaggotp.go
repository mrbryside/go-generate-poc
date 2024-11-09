package handlertp

const SwagBaseTemplate = `
// @Router       #path# [#method#]
// @Accept       json
// @Produce      json`

const SwagSummaryTemplate = `
// @Summary      #content#`

const SwagDescriptionTemplate = `
// @Description      #content#`

const SwagRequestBodyTemplate = `
// @Param request body dto.#request# true "request body"`

const SwagSuccessTemplate = `
// @Success      #statusCode#  {#type#}  dto.#response#`

const SwagFailureTemplate = `
// @Failure      #statusCode#  {#type#}  dto.#response#`

const SwagTagTemplate = `
// @Tags        #content#`
