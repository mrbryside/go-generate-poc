package handlertp

const BindTemplate = `var req #requestName#
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, #responseName#)
	}
`

//func GenRequestBind(requestName string) string {
//	return strings.Replace(bindTemplate, "#requestName#", requestName, -1)
//}
