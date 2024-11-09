package handlergen

import (
	"github.com/mrbryside/go-generate/internal/myfile"
	"github.com/mrbryside/go-generate/internal/myhttp"
	"github.com/mrbryside/go-generate/internal/mymap"
	"strings"
)

func generateRequest(
	template,
	templateGenerate string,
	handlerName string,
	request *mymap.OrderedMap,
	statusCodes []string,
	isStatusCodeStyle bool,
) (string, string) {
	if request == nil || request.Len() == 0 {
		// remove every request block because it's not have request
		template = myfile.RemoveLine(template, "#requestBind#")
		template = myfile.RemoveLine(template, "#requestValidation#")
		return template, templateGenerate
	}
	if isStatusCodeStyle {
		return generateRequestStatusCodeStyle(template, templateGenerate, handlerName, request, statusCodes)
	}

	return generateRequestNonStatusCodeStyle(template, templateGenerate, handlerName, request)
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
	if response == nil || response.Len() == 0 {
		template = ReplaceResponseBlockNoContentForNonStatusCodeStyle(template)
		return template, templateGenerate
	}
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
