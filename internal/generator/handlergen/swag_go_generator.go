package handlergen

import (
	"github.com/mrbryside/go-generate/internal/generator/template/handlertp"
	"github.com/mrbryside/go-generate/internal/myfile"
	"github.com/mrbryside/go-generate/internal/myhttp"
	"strings"
)

func GenerateTempMainHandler(path string, packageName string) string {
	result := GenerateContentMainHandler(packageName)
	result = myfile.RenamePackageGolangFile(result, GenTempGenerateFolderAndPackageName(path))

	return result
}

func GenerateTempUserHandlerWithSwagGoSyntax(path string, currentContent string, htd HandlerTemplateData) string {
	currentContent = myfile.RenamePackageGolangFile(currentContent, GenTempGenerateFolderAndPackageName(path))
	result := generateSwagGoBaseTemplate(currentContent, htd.Api, htd.Method)
	result = generateSwagGoSummaryTemplate(result, htd)
	result = generateSwagGoDescriptionTemplate(result, htd)
	result = generateSwagGoTagTemplate(result, htd)
	result = generateSwagGoRequestTemplate(result, htd)
	result = generateSwagGoSuccessResponseTemplate(result, htd)
	result = generateSwagGoFailureResponseTemplate(result, htd)

	// clear all #swaggo# text
	result = strings.Replace(result, "#swaggo#", "", -1)

	return result
}

func generateSwagGoTagTemplate(currentContent string, htd HandlerTemplateData) string {
	if htd.Tag == "" {
		return currentContent
	}
	template := handlertp.SwagTagTemplate
	template = strings.Replace(template, "#content#", htd.Tag, -1)
	template = AddSwaggoReplaceText(template)
	currentContent = strings.Replace(currentContent, "#swaggo#", template, 1)

	return currentContent
}

func generateSwagGoDescriptionTemplate(currentContent string, htd HandlerTemplateData) string {
	if htd.Description == "" {
		return currentContent
	}
	template := handlertp.SwagDescriptionTemplate
	template = strings.Replace(template, "#content#", htd.Description, -1)
	template = AddSwaggoReplaceText(template)
	currentContent = strings.Replace(currentContent, "#swaggo#", template, 1)

	return currentContent
}

func generateSwagGoSummaryTemplate(currentContent string, htd HandlerTemplateData) string {
	if htd.Summary == "" {
		return currentContent
	}
	template := handlertp.SwagSummaryTemplate
	template = strings.Replace(template, "#content#", htd.Summary, -1)
	template = AddSwaggoReplaceText(template)
	currentContent = strings.Replace(currentContent, "#swaggo#", template, 1)

	return currentContent
}

func generateSwagGoSuccessResponseTemplate(currentContent string, htd HandlerTemplateData) string {
	if htd.Response == nil || htd.Response.Len() == 0 {
		return currentContent
	}
	template := handlertp.SwagSuccessTemplate
	if isStatusCodeStyle(htd.Response) {
		iter := htd.Response.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			if pair.Key == "200" || pair.Key == "201" {
				typeName := htd.Name + myhttp.StatusCodeMap[pair.Key] + "Response"
				template = strings.Replace(template, "#statusCode#", pair.Key, -1)
				template = strings.Replace(template, "#type#", "object", -1)
				template = strings.Replace(template, "#response#", typeName, -1)
				template = AddSwaggoReplaceText(template)
				break
			}
		}
		currentContent = strings.Replace(currentContent, "#swaggo#", template, 1)
		return currentContent
	}
	template = strings.Replace(template, "#statusCode#", "200", -1)
	template = strings.Replace(template, "#type#", "object", -1)
	template = strings.Replace(template, "#response#", htd.Name+"Response", -1)
	template = AddSwaggoReplaceText(template)
	currentContent = strings.Replace(currentContent, "#swaggo#", template, 1)

	return currentContent
}

func generateSwagGoFailureResponseTemplate(currentContent string, htd HandlerTemplateData) string {
	if htd.Response == nil {
		return currentContent
	}
	currentTemplate := ""
	if isStatusCodeStyle(htd.Response) {
		iter := htd.Response.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			if pair.Key >= "400" && pair.Key <= "599" {
				template := handlertp.SwagFailureTemplate
				typeName := htd.Name + myhttp.StatusCodeMap[pair.Key] + "Response"
				template = strings.Replace(template, "#statusCode#", pair.Key, -1)
				template = strings.Replace(template, "#type#", "object", -1)
				template = strings.Replace(template, "#response#", typeName, -1)
				template = AddSwaggoReplaceText(template)
				currentTemplate += template
			}
		}
		currentContent = strings.Replace(currentContent, "#swaggo#", currentTemplate, 1)
		return currentContent
	}
	return currentContent
}

func generateSwagGoBaseTemplate(currentContent, apiPath, apiMethod string) string {
	template := handlertp.SwagBaseTemplate
	template = strings.Replace(template, "#path#", apiPath, -1)
	template = strings.Replace(template, "#method#", apiMethod, -1)
	template = AddSwaggoReplaceText(template)
	currentContent = strings.Replace(currentContent, "#swaggo#", template, 1)

	return currentContent
}

func generateSwagGoRequestTemplate(currentContent string, htd HandlerTemplateData) string {
	if htd.Request == nil || htd.Request.Len() == 0 {
		return currentContent
	}
	template := handlertp.SwagRequestBodyTemplate
	template = strings.Replace(template, "#request#", htd.Name+"Request", -1)
	template = AddSwaggoReplaceText(template)
	currentContent = strings.Replace(currentContent, "#swaggo#", template, 1)

	return currentContent
}

func AddSwaggoReplaceText(content string) string {
	return content + "#swaggo#"
}

func ReplaceSwaggoText(content string) string {
	return strings.Replace(content, "#swaggo#", "", -1)
}
