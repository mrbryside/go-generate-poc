package handlergen

import "github.com/mrbryside/go-generate/internal/utils/mystr"

func GenHandlerFileNameFromHandlerTemplate(htd HandlerTemplateData) string {
	return mystr.ToSnakeCase(htd.Name)
}

func GenHandlerFunctionExportedNameFromHandlerTemplate(htd HandlerTemplateData) string {
	return mystr.CapitalizeFirstLetter(htd.Name)
}
