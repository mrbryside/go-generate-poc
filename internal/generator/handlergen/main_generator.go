package handlergen

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mrbryside/go-generate/internal/generator/template/handlertp"
	"github.com/mrbryside/go-generate/internal/utils/myfile"
	"github.com/mrbryside/go-generate/internal/utils/mystr"
)

type HandlerTemplatedDataError struct {
	HandlerNameTemplateData HandlerTemplateData
	Errors                  []error
}

type MandaToryError struct {
	Path  string
	Error error
}

func MainGenerateHandler(path string) (report Report) {
	basePathToGenerate := GenBasePath(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return Report{BasePathOfJsonSpec: basePathToGenerate, MandaToryError: MandaToryError{Path: basePathToGenerate, Error: err}}
	}

	var handlerTemplates []HandlerTemplateData
	packageName := filepath.Base(filepath.Dir(path))
	err = json.Unmarshal(data, &handlerTemplates)
	if err != nil {
		return Report{BasePathOfJsonSpec: basePathToGenerate, MandaToryError: MandaToryError{Path: basePathToGenerate, Error: err}}
	}

	var pathToGenerateErrors []PathToGenerateError
	var handlerTemplatesSuccessForGenerateRoute []HandlerTemplateData
	var contentToGenerateTemps []string
	var pathToGenerateTemps []string
	var pathToGenerates []string
	var contentToGenerates []string
	var handlerTemplatedSuccessRoutes []HandlerTemplateData
	var handlerTemplatedFailedRoutes []HandlerTemplatedDataError
	for _, ht := range handlerTemplates {
		err := ValidateHandler(ht)
		if err != nil {
			// need to get path and file that we generate and add to knab_logs.json because we dont wanna delete user file that write wrong format of spec
			// pathToGenerateErrors will use by function that call MainGenerateHandler
			fileUserHandlerName := GenGoFileNameInBasePath(path, GenHandlerFileNameFromHandlerTemplate(ht))
			fileGeneratedHandlerName := GenGoFileNameGeneratedInDtoBasePath(path, GenHandlerFileNameFromHandlerTemplate(ht))
			pathToGenerateErrors = append(pathToGenerateErrors, PathToGenerateError{Path: fileUserHandlerName, Error: err})
			pathToGenerateErrors = append(pathToGenerateErrors, PathToGenerateError{Path: fileGeneratedHandlerName, Error: err})

			// add to failed routes
			handlerTemplatedFailedRoutes = append(handlerTemplatedFailedRoutes, HandlerTemplatedDataError{
				ht,
				[]error{},
			})
			continue
		}
		handlerTemplatesSuccessForGenerateRoute = append(handlerTemplatesSuccessForGenerateRoute, ht)
		// generate section
		userHandlerContent, generatedHandlerContent := GenerateContentBothUserHandlerAndGeneratedHandler(packageName, ht) // generate handler
		// TODO: generate query params, path params

		pathToGenerates, contentToGenerates = AddContentForUserHandlerAndGeneratedHandler(path, pathToGenerates, contentToGenerates, userHandlerContent, generatedHandlerContent, ht)
		pathToGenerateTemps, contentToGenerateTemps = AddContentForTempUserHandlerAndGeneratedHandler(path, pathToGenerateTemps, contentToGenerateTemps, userHandlerContent, generatedHandlerContent, ht)

		// add to success routes
		handlerTemplatedSuccessRoutes = append(handlerTemplatedSuccessRoutes, ht)
	}
	pathToGenerates, contentToGenerates = AddContentForMainHandlerAndRouteFile(path, packageName, pathToGenerates, contentToGenerates, handlerTemplatesSuccessForGenerateRoute)
	pathToGenerateTemps, contentToGenerateTemps = AddContentForTempMainHandlerAndRouteFile(path, packageName, pathToGenerateTemps, contentToGenerateTemps, handlerTemplatesSuccessForGenerateRoute)

	pathToGenerates, contentToGenerates = AddContentForValidatorHelper(path, packageName, pathToGenerates, contentToGenerates)
	pathToGenerateTemps, contentToGenerateTemps = AddContentForTempValidatorHelper(path, pathToGenerateTemps, contentToGenerateTemps)

	// write all temp output paths
	indexGenerateErrors := make([]int, 0)
	for i, pgt := range pathToGenerateTemps {
		err = os.WriteFile(pgt, []byte(contentToGenerateTemps[i]), 0644)
		if err != nil {
			pathToGenerateErrors = append(pathToGenerateErrors, PathToGenerateError{Path: pgt, Error: err})
			indexGenerateErrors = append(indexGenerateErrors, i)
			continue
		}

		cmd := exec.Command("go", "fmt", pgt)
		err = cmd.Run()
		if err != nil {
			pathToGenerateErrors = append(pathToGenerateErrors, PathToGenerateError{Path: pgt, Error: errors.New("go format error")})
			indexGenerateErrors = append(indexGenerateErrors, i)
			continue
		}

		cmd = exec.Command("goimports", "-w", pgt)
		err = cmd.Run()
		if err != nil {
			pathToGenerateErrors = append(pathToGenerateErrors, PathToGenerateError{Path: pgt, Error: errors.New("go imports error")})
			indexGenerateErrors = append(indexGenerateErrors, i)
			continue
		}
	}

	// run go vet to verify not have any error in folder temp
	pathToRunVetInTempFolder := "./" + basePathToGenerate + "/" + GenTempGenerateFolderAndPackageName(path)
	cmd := exec.Command("go", "vet", pathToRunVetInTempFolder)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Report{BasePathOfJsonSpec: basePathToGenerate, MandaToryError: MandaToryError{Path: basePathToGenerate, Error: errors.New(fmt.Sprintf("%v: %s", err, string(output)))}}
	}

	// filter out failed from handlerTemplatedSuccessRoutes matching with pathToGenerateErrors
	var handlerTemplatedSuccessRoutesFilterOutFailed []HandlerTemplateData
	for _, ht := range handlerTemplatedSuccessRoutes {
		match := false
		for _, pgErr := range pathToGenerateErrors {
			if strings.Contains(pgErr.Path, mystr.ToSnakeCase(ht.Name)) {
				match = true
				handlerTemplatedFailedRoutes = append(handlerTemplatedFailedRoutes, HandlerTemplatedDataError{
					HandlerNameTemplateData: ht,
				})
				continue
			}
		}
		if !match {
			handlerTemplatedSuccessRoutesFilterOutFailed = append(handlerTemplatedSuccessRoutesFilterOutFailed, ht)
		}
	}

	// add more error from pathToGenerateErrors, handlerTemplatedFailedRoutes to handlerTemplateRoutesAddMoreFailedFromGenerate
	var handlerTemplateRoutesAddMoreFailedFromGenerate []HandlerTemplatedDataError
	mappedHandlerTemplateRoutesAddMoreFailedFromGenerate := make(map[string]HandlerTemplatedDataError)
	for _, ht := range handlerTemplatedFailedRoutes {
		match := false
		for _, pgErr := range pathToGenerateErrors {
			if strings.Contains(pgErr.Path, mystr.ToSnakeCase(ht.HandlerNameTemplateData.Name)) {
				match = true
				newHtAddedErr := HandlerTemplatedDataError{
					HandlerNameTemplateData: ht.HandlerNameTemplateData,
					Errors:                  append(ht.Errors, pgErr.Error),
				}
				mappedHandlerTemplateRoutesAddMoreFailedFromGenerate[ht.HandlerNameTemplateData.Name] = newHtAddedErr
			}
		}
		if !match {
			mappedHandlerTemplateRoutesAddMoreFailedFromGenerate[ht.HandlerNameTemplateData.Name] = ht
		}
	}
	for _, v := range mappedHandlerTemplateRoutesAddMoreFailedFromGenerate {
		handlerTemplateRoutesAddMoreFailedFromGenerate = append(handlerTemplateRoutesAddMoreFailedFromGenerate, v)
	}

	// construct pathToGenerateWithoutError
	pathToGenerateWithoutError := make([]string, 0)
	for i, ptg := range pathToGenerates {
		match := false
		for _, idxErr := range indexGenerateErrors {
			if idxErr == i {
				match = true
				continue
			}
		}
		if !match {
			pathToGenerateWithoutError = append(pathToGenerateWithoutError, ptg)
		}
	}

	// write all pathToGenerates with no error
	for i, ptg := range pathToGenerates {
		// check if it's a routes_gen.go (it's call function of handler that can be error inside that function if have any error skip to gen it
		if strings.Contains(ptg, RouteFileName+".go") {
			if len(indexGenerateErrors) > 0 {
				continue
			}
		}
		// check exist in error index from generate temp? if yes continue
		foundErr := false
		for _, idx := range indexGenerateErrors {
			if idx == i {
				foundErr = true
				break
			}
		}
		if foundErr {
			continue
		}
		_ = os.WriteFile(ptg, []byte(contentToGenerates[i]), 0644)

		cmd := exec.Command("go", "fmt", ptg)
		_ = cmd.Run()

		cmd = exec.Command("goimports", "-w", ptg)
		_ = cmd.Run()
	}

	// run swag init
	isSwagInitSuccess := true
	swagInitError := error(nil)
	err = runSwagGoInit(path)
	if err != nil {
		swagInitError = err
		isSwagInitSuccess = false
	}

	return Report{
		PathToGenerateSuccess:       pathToGenerateWithoutError,
		PathToGenerateError:         pathToGenerateErrors,
		HandlerTemplateSuccessRoute: handlerTemplatedSuccessRoutesFilterOutFailed,
		HandlerTemplateErrorRoute:   handlerTemplateRoutesAddMoreFailedFromGenerate,
		BasePathOfJsonSpec:          GenBasePath(path),
		SwagGenerateReport: SwagGenerateReport{
			isSuccess: isSwagInitSuccess,
			Error:     swagInitError,
		},
		MandaToryError: MandaToryError{Error: nil},
	}
}

func AddContentForTempValidatorHelper(path string, pathToGenerates []string, contentToGenerates []string) ([]string, []string) {
	packageName := GenTempGenerateFolderAndPackageName(path)
	validatorHelperPath := GenGoFileNameInTempBasePath(path, ValidatorFileName)
	validationHelperContent := handlertp.ValidationHelperTemplate
	validationHelperContent = myfile.RenamePackageGolangFileContent(validationHelperContent, packageName)

	pathToGenerates = append(pathToGenerates, validatorHelperPath)
	contentToGenerates = append(contentToGenerates, validationHelperContent)

	return pathToGenerates, contentToGenerates
}

func AddContentForValidatorHelper(path string, packageName string, pathToGenerateTemps []string, contentToGenerateTemps []string) ([]string, []string) {
	validatorHelperPath := GenGoFileNameInBasePath(path, ValidatorFileName)
	validationHelperContent := handlertp.ValidationHelperTemplate
	validationHelperContent = myfile.RenamePackageGolangFileContent(validationHelperContent, packageName)

	pathToGenerateTemps = append(pathToGenerateTemps, validatorHelperPath)
	contentToGenerateTemps = append(contentToGenerateTemps, validationHelperContent)

	return pathToGenerateTemps, contentToGenerateTemps
}

func AddContentForTempUserHandlerAndGeneratedHandler(path string, pathToGenerateTemps []string, contentToGenerateTemps []string, userHandlerContent string, generatedHandlerContent string, ht HandlerTemplateData) ([]string, []string) {
	// generate user handler for temp
	userHandlerTempContent := GenerateTempUserHandlerWithSwagGoSyntax(path, userHandlerContent, ht)
	userHandlerTempContent = strings.Replace(userHandlerTempContent, handlertp.ModuleReplaceName, GenModuleDtoTempPath(path), -1)

	userHandlerTempPath := GenGoFileNameInTempBasePath(path, GenHandlerFileNameFromHandlerTemplate(ht))
	pathToGenerateTemps = append(pathToGenerateTemps, userHandlerTempPath)
	contentToGenerateTemps = append(contentToGenerateTemps, userHandlerTempContent)

	// generate generated handler for temp
	generateHandlerTempPath := GenGoFileNameGeneratedInDtoTempBasePath(path, GenHandlerFileNameFromHandlerTemplate(ht))
	pathToGenerateTemps = append(pathToGenerateTemps, generateHandlerTempPath)
	contentToGenerateTemps = append(contentToGenerateTemps, generatedHandlerContent)

	return pathToGenerateTemps, contentToGenerateTemps
}

func AddContentForUserHandlerAndGeneratedHandler(path string, pathToGenerates []string, contentToGenerates []string, userHandlerContent string, generatedHandlerContent string, ht HandlerTemplateData) ([]string, []string) {
	// generate user handler
	userHandlerContent = strings.Replace(userHandlerContent, handlertp.SwagGoReplaceName, "", -1)
	userHandlerContent = strings.Replace(userHandlerContent, handlertp.ModuleReplaceName, GenModuleDtoPath(path), -1)

	userHandlerPath := GenGoFileNameInBasePath(path, GenHandlerFileNameFromHandlerTemplate(ht))
	pathToGenerates = append(pathToGenerates, userHandlerPath)
	contentToGenerates = append(contentToGenerates, userHandlerContent)

	// generate generated handler
	generateHandlerPath := GenGoFileNameGeneratedInDtoBasePath(path, GenHandlerFileNameFromHandlerTemplate(ht))
	pathToGenerates = append(pathToGenerates, generateHandlerPath)
	contentToGenerates = append(contentToGenerates, generatedHandlerContent)

	return pathToGenerates, contentToGenerates
}

func AddContentForMainHandlerAndRouteFile(path string, packageName string, pathToGenerates []string, contentToGenerates []string, htds []HandlerTemplateData) ([]string, []string) {
	// generate main handler file
	mainHandlerContent := GenerateContentMainHandler(packageName)
	mainHandlerOutputPath := GenGoFileNameInBasePath(path, HandlerFileName)
	pathToGenerates = append(pathToGenerates, mainHandlerOutputPath)
	contentToGenerates = append(contentToGenerates, mainHandlerContent)

	// generate handler routes file
	routeContents := GenerateContentRoutes(htds, packageName)
	routeOutputPath := GenGoFileNameInBasePath(path, RouteFileName)
	pathToGenerates = append(pathToGenerates, routeOutputPath)
	contentToGenerates = append(contentToGenerates, routeContents)

	return pathToGenerates, contentToGenerates
}

func AddContentForTempMainHandlerAndRouteFile(path string, packageName string, pathToGenerateTemps []string, contentToGenerateTemps []string, htds []HandlerTemplateData) ([]string, []string) {
	// generate main handler for temp generated
	mainHandlerTempPath := GenGoFileNameInTempBasePath(path, HandlerFileName)
	mainHandlerTempContent := GenerateTempMainHandlerContent(path, packageName)
	mainHandlerTempContent = myfile.RenamePackageGolangFileContent(mainHandlerTempContent, GenTempGenerateFolderAndPackageName(path))
	pathToGenerateTemps = append(pathToGenerateTemps, mainHandlerTempPath)
	contentToGenerateTemps = append(contentToGenerateTemps, mainHandlerTempContent)

	// generate handler routes file
	routeContents := GenerateContentRoutes(htds, packageName)
	routeOutputPath := GenGoFileNameInTempBasePath(path, RouteFileName)
	routeContents = myfile.RenamePackageGolangFileContent(routeContents, GenTempGenerateFolderAndPackageName(path))
	pathToGenerateTemps = append(pathToGenerateTemps, routeOutputPath)
	contentToGenerateTemps = append(contentToGenerateTemps, routeContents)

	return pathToGenerateTemps, contentToGenerateTemps
}

func GenerateContentBothUserHandlerAndGeneratedHandler(packageName string, htd HandlerTemplateData) (string, string) {
	template := handlertp.Template

	template = strings.Replace(template, handlertp.PackageNameReplaceName, packageName, -1)
	template = strings.Replace(template, handlertp.HandlerFuncNameReplaceName, GenHandlerFunctionExportedNameFromHandlerTemplate(htd), -1)
	template = strings.Replace(template, handlertp.ResponseNameReplaceName, GetHandlerRequestName(htd.Name), -1)
	template = strings.Replace(template, handlertp.RequestNameReplaceName, GetHandlerResponseName(htd.Name), -1)

	templateGenerate := handlertp.TemplateGenerate
	templateGenerate = strings.Replace(templateGenerate, handlertp.PackageNameReplaceName, handlertp.DtoFolderAndPackageName, -1)
	templateGenerate = strings.Replace(templateGenerate, handlertp.HandlerFuncNameReplaceName, GenHandlerFunctionExportedNameFromHandlerTemplate(htd), -1)

	response, statusCodes, isStatusCodeStyle := transformStatusCodeStyleResponse(htd.Response)
	template, templateGenerate = generateRequest(template, templateGenerate, htd.Name, htd.Request, statusCodes, isStatusCodeStyle)
	template, templateGenerate = generateResponse(template, templateGenerate, htd.Name, response, statusCodes, isStatusCodeStyle)

	if htd.Request == nil || htd.Request.Len() == 0 {
		return template, templateGenerate
	}

	return template, templateGenerate
}

func GenerateContentMainHandler(packageName string) string {
	template := handlertp.HandlerMainTemplate
	template = strings.Replace(template, handlertp.PackageNameReplaceName, packageName, -1)
	return template
}

func GenerateContentRoutes(htds []HandlerTemplateData, packageName string) string {
	var routeContents string
	for idx, htd := range htds {
		echoRouteTemplate := handlertp.EchoRouteTemplate
		echoRouteTemplate = strings.Replace(echoRouteTemplate, handlertp.HandlerMethodReplaceName, strings.ToUpper(htd.Method), -1)
		echoRouteTemplate = strings.Replace(echoRouteTemplate, handlertp.HandlerRouteReplaceName, htd.Api, -1)
		echoRouteTemplate = strings.Replace(echoRouteTemplate, handlertp.HandlerFuncNameReplaceName, mystr.CapitalizeFirstLetter(htd.Name), -1)
		routeContents += echoRouteTemplate
		if idx != len(htds)-1 {
			routeContents += "\n"
		}
	}
	routeTemplate := handlertp.RouteTemplate
	routeTemplate = strings.Replace(routeTemplate, handlertp.PackageNameReplaceName, packageName, -1)
	return strings.Replace(routeTemplate, handlertp.RouteReplaceName, routeContents, -1)
}
