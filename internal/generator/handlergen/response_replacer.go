package handlergen

import (
	"github.com/mrbryside/go-generate/internal/generator/template/handlertp"
	"github.com/mrbryside/go-generate/internal/myhttp"
	"strings"
)

func ReplaceSuccessResponseBlockForStatusCodeStyle(currentContent string, statusCode []string, handlerName string) string {
	successTemplate := ""
	for _, code := range statusCode {
		statusCodeMap := myhttp.StatusCodeMap[code]
		if code == "200" {
			successTemplate = handlertp.SuccessJsonResponseTemplate
			successTemplate = strings.Replace(successTemplate, "#httpStatus#", "http.StatusOK", -1)
			successTemplate = strings.Replace(successTemplate, "#responseName#", handlerName+statusCodeMap+"Response{}", -1)
			break
		}
		if code == "201" {
			successTemplate = handlertp.SuccessJsonResponseTemplate
			successTemplate = strings.Replace(successTemplate, "#httpStatus#", "http.StatusCreated", -1)
			successTemplate = strings.Replace(successTemplate, "#responseName#", handlerName+statusCodeMap+"Response{}", -1)
			break
		}
	}
	// validateTemplate is empty means not found status code 200 using default validation template
	if successTemplate == "" {
		return ReplaceResponseBlockNoContentForNonStatusCodeStyle(currentContent)
	}
	currentContent = strings.Replace(currentContent, "#returnSuccess#", successTemplate, -1)
	return currentContent
}

func ReplaceResponseBlockNoContentForNonStatusCodeStyle(currentContent string) string {
	successNoContentTemplate := handlertp.SuccessNoContentResponseTemplate
	successNoContentTemplate = strings.Replace(successNoContentTemplate, "#httpStatus#", "http.StatusOK", -1)
	currentContent = strings.Replace(currentContent, "#returnSuccess#", successNoContentTemplate, -1)
	return currentContent
}

func ReplaceSuccessResponseBlockForNonStatusCodeStyle(currentContent string, handlerName string) string {
	successTemplate := handlertp.SuccessJsonResponseTemplate
	successTemplate = strings.Replace(successTemplate, "#httpStatus#", "http.StatusOK", -1)
	successTemplate = strings.Replace(successTemplate, "#responseName#", handlerName+"Response{}", -1)

	currentContent = strings.Replace(currentContent, "#returnSuccess#", successTemplate, -1)
	return currentContent
}

func ReplaceValidationBlockForStatusCodeStyle(currentContent string, statusCodes []string, handlerName string) string {
	validationTemplate := ""
	for _, code := range statusCodes {
		statusCodeMap := myhttp.StatusCodeMap[code]
		if code == "400" {
			validationTemplate = handlertp.RequestValidationForStatusCodeStyleTemplate
			validationTemplate = strings.Replace(validationTemplate, "#handlerFuncName#", handlerName, -1)
			validationTemplate = strings.Replace(validationTemplate, "#responseName#", handlerName+statusCodeMap+"Response{}", -1)
		}
	}
	// validateTemplate is empty means not found status code 400 using default validation template
	if validationTemplate == "" {
		return ReplaceValidationBlockForNonStatusCodeStyle(currentContent, handlerName)
	}

	currentContent = strings.Replace(currentContent, "#requestValidation#", validationTemplate, -1)
	return currentContent
}

func ReplaceValidationBlockForNonStatusCodeStyle(currentContent string, handlerName string) string {
	validationTemplate := ""
	validationTemplate = handlertp.RequestValidationTemplate
	validationTemplate = strings.Replace(validationTemplate, "#handlerFuncName#", handlerName, -1)

	currentContent = strings.Replace(currentContent, "#requestValidation#", validationTemplate, -1)
	return currentContent
}

func ReplaceBindingBlockForStatusCodeStyle(currentContent string, statusCodes []string, handlerName string) string {
	bindingTemplate := ""
	for _, code := range statusCodes {
		statusCodeMap := myhttp.StatusCodeMap[code]
		if code == "400" {
			bindingTemplate = handlertp.BindTemplate
			bindingTemplate = strings.Replace(bindingTemplate, "#requestName#", handlerName+"Request", -1)
			bindingTemplate = strings.Replace(bindingTemplate, "#responseName#", handlerName+statusCodeMap+"Response{}", -1)
		}
	}
	// validateTemplate is empty means not found status code 400 using default validation template
	if bindingTemplate == "" {
		return ReplaceBindingBlockForNonStatusCodeStyle(currentContent, handlerName)
	}

	currentContent = strings.Replace(currentContent, "#requestBind#", bindingTemplate, -1)
	return currentContent
}

func ReplaceBindingBlockForNonStatusCodeStyle(currentContent string, handlerName string) string {
	bindingTemplate := handlertp.BindTemplate
	bindingTemplate = strings.Replace(bindingTemplate, "#requestName#", handlerName+"Request", -1)
	bindingTemplate = strings.Replace(bindingTemplate, "#responseName#", "err", -1)

	currentContent = strings.Replace(currentContent, "#requestBind#", bindingTemplate, -1)
	return currentContent
}
