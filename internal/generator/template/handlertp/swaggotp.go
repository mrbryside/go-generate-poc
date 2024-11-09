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
// @Param request body #request# true "request body"`

const SwagSuccessTemplate = `
// @Success      #statusCode#  {#type#}  #response#`

const SwagFailureTemplate = `
// @Failure      #statusCode#  {#type#}  #response#`

const SwagTagTemplate = `
// @Tags        #content#`
