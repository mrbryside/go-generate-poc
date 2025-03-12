package handlertp

import "fmt"

var SuccessJsonResponseTemplate = fmt.Sprintf(`return ctx.JSON(%s, %s.%s)`, HttpStatusReplaceName, DtoFolderAndPackageName, ResponseNameReplaceName)
var SuccessNoContentResponseTemplate = fmt.Sprintf(`return ctx.NoContent(%s)`, HttpStatusReplaceName)
