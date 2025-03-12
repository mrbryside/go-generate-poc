package handlertp

import "fmt"

const ()

var SwagBaseTemplate = fmt.Sprintf(`
// @Router       %s [%s]
// @Accept       json
// @Produce      json`, SwagGoPathReplaceName, SwagGoMethodReplaceName)

var SwagSummaryTemplate = fmt.Sprintf(`
// @Summary      %s`, SwagGoContentReplaceName)

var SwagDescriptionTemplate = fmt.Sprintf(`
// @Description      %s`, SwagGoContentReplaceName)

var SwagRequestBodyTemplate = fmt.Sprintf(`
// @Param request body %s.%s true "request body"`, DtoFolderAndPackageName, SwagGoRequestReplaceName)

var SwagSuccessTemplate = fmt.Sprintf(`
// @Success      %s  {%s}  %s.%s`,
	SwagGoStatusCodeReplaceName,
	SwagGoTypeReplaceName,
	DtoFolderAndPackageName,
	SwagGoResponseReplaceName,
)

var SwagFailureTemplate = fmt.Sprintf(`
// @Failure      %s  {%s}  %s.%s`,
	SwagGoStatusCodeReplaceName,
	SwagGoTypeReplaceName,
	DtoFolderAndPackageName,
	SwagGoResponseReplaceName,
)

var SwagTagTemplate = fmt.Sprintf(`
// @Tags        %s`, SwagGoContentReplaceName)
