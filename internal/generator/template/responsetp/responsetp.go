package responsetp

import (
	"strings"
)

const responseTp = `return ctx.JSON(#status#, #response#)`
const responseTpOnlyStatus = `return ctx.NoContent(#status#)`

func GenResponseTp(httpStatus string, response string) string {
	result := strings.Replace(responseTp, "#status#", httpStatus, 1)
	return strings.Replace(result, "#response#", response+"{}", 1)
}

func GenResponseTpOnlyStatus(httpStatus string) string {
	return strings.Replace(responseTpOnlyStatus, "#status#", httpStatus, 1)
}
