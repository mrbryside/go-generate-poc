package handlergen

import (
	"github.com/mrbryside/go-generate/internal/myfile"
	"github.com/mrbryside/go-generate/internal/myhttp"
	"github.com/mrbryside/go-generate/internal/mymap"
	"strings"

	"github.com/mrbryside/go-generate/internal/generator/template/handlertp"
	"github.com/mrbryside/go-generate/internal/mystr"
)

func GenerateHandler(packageName string, htd HandlerTemplateData) (string, string) {
	template := handlertp.Template
	template = strings.Replace(template, "handlertp", packageName, -1)
	template = strings.Replace(template, "#handlerName#", mystr.CapitalizeFirstLetter(htd.Name)+"Handler", -1)
	template = strings.Replace(template, "#handlerFuncName#", mystr.CapitalizeFirstLetter(htd.Name), -1)
	template = strings.Replace(template, "#requestName#", htd.Name+"Request", -1)
	template = strings.Replace(template, "#responseName#", htd.Name+"Response", -1)

	templateGenerate := handlertp.TemplateGenerate
	templateGenerate = strings.Replace(templateGenerate, "handlertp", packageName, -1)
	templateGenerate = strings.Replace(templateGenerate, "#handlerFuncName#", mystr.CapitalizeFirstLetter(htd.Name), -1)

	response, statusCodes, isStatusCodeStyle := transformStatusCodeStyleResponse(htd.Response)
	template, templateGenerate = generateRequest(template, templateGenerate, htd.Name, htd.Request, statusCodes, isStatusCodeStyle)
	template, templateGenerate = generateResponse(template, templateGenerate, htd.Name, response, statusCodes, isStatusCodeStyle)

	// add some content to last line after generate struct to template_generate
	if htd.Request != nil {
		// add validation helper
		validationHelperContent := strings.Replace(handlertp.ValidationHelperTemplate, "#handlerFuncName#", htd.Name, -1)
		templateGenerate = myfile.AddContentToLastLine(templateGenerate, validationHelperContent)
	}

	return template, templateGenerate
}

func GenerateMainHandler() string {
	return handlertp.HandlerMainTemplate
}

func GenerateHandlerRoutes(htds []HandlerTemplateData) string {
	var routeContents string
	for idx, htd := range htds {
		echoRouteTemplate := handlertp.EchoRouteTemplate
		echoRouteTemplate = strings.Replace(echoRouteTemplate, "#handlerMethod#", strings.ToUpper(htd.Method), -1)
		echoRouteTemplate = strings.Replace(echoRouteTemplate, "#handlerRoute#", htd.Api, -1)
		echoRouteTemplate = strings.Replace(echoRouteTemplate, "#handlerFuncName#", mystr.CapitalizeFirstLetter(htd.Name), -1)
		routeContents += echoRouteTemplate
		if idx != len(htds)-1 {
			routeContents += "\n"
		}
	}
	routeTemplate := handlertp.RouteTemplate
	return strings.Replace(routeTemplate, "#route#", routeContents, -1)
}

func generateRequest(
	template,
	templateGenerate string,
	handlerName string,
	request *mymap.OrderedMap,
	statusCodes []string,
	isStatusCodeStyle bool,
) (string, string) {
	if request != nil && request.Len() != 0 {
		if isStatusCodeStyle {
			return generateRequestStatusCodeStyle(template, templateGenerate, handlerName, request, statusCodes)
		}
		return generateRequestNonStatusCodeStyle(template, templateGenerate, handlerName, request)
	}
	// remove every request block because it's not have request
	template = myfile.RemoveLine(template, "#requestBind#")
	template = myfile.RemoveLine(template, "#requestValidation#")

	return template, templateGenerate
}

func generateRequestStatusCodeStyle(
	template string,
	templateGenerate string,
	handlerName string,
	request *mymap.OrderedMap,
	statusCodes []string,
) (string, string) {
	// generate struct for request
	fieldsRequestString, newRequestStructs := generateStructFields("Request", handlerName, request, []string{}, []myfile.NewStruct{}, "")
	templateGenerate = strings.Replace(templateGenerate, "#requestFields#", fieldsRequestString, -1)
	templateGenerate = myfile.AddStructToLastLine(templateGenerate, fieldsRequestString, handlerName+"Request")
	if len(newRequestStructs) != 0 {
		for _, newRequestStruct := range newRequestStructs {
			templateGenerate = myfile.AddStructToLastLine(templateGenerate, newRequestStruct.Fields, newRequestStruct.Name)
		}
	}

	// generate validation and binding block
	template = ReplaceValidationBlockForStatusCodeStyle(template, statusCodes, handlerName)
	template = ReplaceBindingBlockForStatusCodeStyle(template, statusCodes, handlerName)
	return template, templateGenerate
}

func generateRequestNonStatusCodeStyle(
	template string,
	templateGenerate string,
	handlerName string,
	request *mymap.OrderedMap,
) (string, string) {
	// generate struct for request
	fieldsRequestString, newRequestStructs := generateStructFields("Request", handlerName, request, []string{}, []myfile.NewStruct{}, "")
	templateGenerate = strings.Replace(templateGenerate, "#requestFields#", fieldsRequestString, -1)
	templateGenerate = myfile.AddStructToLastLine(templateGenerate, fieldsRequestString, handlerName+"Request")
	if len(newRequestStructs) != 0 {
		for _, newRequestStruct := range newRequestStructs {
			templateGenerate = myfile.AddStructToLastLine(templateGenerate, newRequestStruct.Fields, newRequestStruct.Name)
		}
	}

	// generate validation and binding block
	template = ReplaceValidationBlockForNonStatusCodeStyle(template, handlerName)
	template = ReplaceBindingBlockForNonStatusCodeStyle(template, handlerName)
	return template, templateGenerate
}

func generateResponse(
	template,
	templateGenerate,
	handlerName string,
	response *mymap.OrderedMap,
	statusCodes []string,
	isStatusCodeStyle bool,
) (string, string) {
	if isStatusCodeStyle {
		return generateResponseStatusCodeStyle(template, templateGenerate, handlerName, response, statusCodes)
	}
	return generateResponseNonStatusCodeStyle(template, templateGenerate, handlerName, response)
}

func generateResponseStatusCodeStyle(
	template string,
	templateGenerate string,
	handlerName string,
	response *mymap.OrderedMap,
	statusCodes []string,
) (string, string) {
	if response == nil {
		template = ReplaceResponseBlockNoContentForNonStatusCodeStyle(template)
		return template, templateGenerate
	}

	_, newResponseStructs := generateStructFields("Response", handlerName, response, []string{}, []myfile.NewStruct{}, "")
	if len(newResponseStructs) != 0 {
		for _, newResponseStruct := range newResponseStructs {
			templateGenerate = myfile.AddStructToLastLine(templateGenerate, newResponseStruct.Fields, newResponseStruct.Name)
		}
	}
	template = ReplaceSuccessResponseBlockForStatusCodeStyle(template, statusCodes, handlerName)

	return template, templateGenerate
}

func generateResponseNonStatusCodeStyle(
	template,
	templateGenerate,
	handlerName string,
	response *mymap.OrderedMap,
) (string, string) {
	if response == nil {
		template = ReplaceResponseBlockNoContentForNonStatusCodeStyle(template)
		return template, templateGenerate
	}

	fieldsResponseString, newResponseStructs := generateStructFields("Response", handlerName, response, []string{}, []myfile.NewStruct{}, "")
	if len(newResponseStructs) != 0 {
		for _, newResponseStruct := range newResponseStructs {
			templateGenerate = myfile.AddStructToLastLine(templateGenerate, newResponseStruct.Fields, newResponseStruct.Name)
		}
	}
	// only create `handlerName+Response` struct when not in status code style
	// status code style will have own response struct for specific status code
	templateGenerate = myfile.AddStructToLastLine(templateGenerate, fieldsResponseString, handlerName+"Response")
	template = ReplaceSuccessResponseBlockForNonStatusCodeStyle(template, handlerName)

	return template, templateGenerate
}

func transformStatusCodeStyleResponse(response *mymap.OrderedMap) (*mymap.OrderedMap, []string, bool) {
	if !isStatusCodeStyle(response) {
		return response, []string{}, false
	}
	var newResponse = mymap.NewOrderedMap()
	var responseCode = make([]string, 0)

	iter := response.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		responseCode = append(responseCode, pair.Key)
		newResponse.Set(myhttp.StatusCodeMap[pair.Key], pair.Value)
	}

	return newResponse, responseCode, true
}
