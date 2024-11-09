package handlergen

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mrbryside/go-generate/internal/mygo"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mrbryside/go-generate/internal/generator/template/handlertp"
	"github.com/mrbryside/go-generate/internal/myfile"
	"github.com/mrbryside/go-generate/internal/mystr"
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
	currentPathFolderToGenerate := strings.TrimSuffix(path, "/handler.json")
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return Report{BasePathOfJsonSpec: currentPathFolderToGenerate, MandaToryError: MandaToryError{Path: currentPathFolderToGenerate, Error: err}}
	}

	var handlerTemplates []HandlerTemplateData
	packageName := filepath.Base(filepath.Dir(path))
	err = json.Unmarshal(data, &handlerTemplates)
	if err != nil {
		return Report{BasePathOfJsonSpec: currentPathFolderToGenerate, MandaToryError: MandaToryError{Path: currentPathFolderToGenerate, Error: err}}
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
			currentPath := strings.TrimSuffix(path, "handler.json")
			fileUserHandlerName := currentPath + mystr.ToSnakeCase(ht.Name) + ".go"
			fileGeneratedHandlerName := currentPath + "dto/" + mystr.ToSnakeCase(ht.Name) + "_gen.go"
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

	// run go vet to verify not have any error the folder temp
	pathToRunVetInTempFolder := "./" + currentPathFolderToGenerate + "/" + GenTempGenerateFolderAndPackageName(path)
	cmd := exec.Command("go", "vet", pathToRunVetInTempFolder)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Report{BasePathOfJsonSpec: currentPathFolderToGenerate, MandaToryError: MandaToryError{Path: currentPathFolderToGenerate, Error: errors.New(fmt.Sprintf("%v: %s", err, string(output)))}}
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

	// add more error from pathToGenerateErrors,handlerTemplatedFailedRoutes to handlerTemplateRoutesAddMoreFailedFromGenerate
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

	//// construct pathToGenerateWithoutError
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
		if strings.Contains(ptg, "routes_gen.go") {
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
	validatorHelperPath := strings.TrimSuffix(path, "handler.json") + fmt.Sprintf("%s/", GenTempGenerateFolderAndPackageName(path)) + "validator.go"
	validationHelperContent := handlertp.ValidationHelperTemplate
	validationHelperContent = myfile.RenamePackageGolangFile(validationHelperContent, packageName)

	pathToGenerates = append(pathToGenerates, validatorHelperPath)
	contentToGenerates = append(contentToGenerates, validationHelperContent)

	return pathToGenerates, contentToGenerates
}

func AddContentForValidatorHelper(path string, packageName string, pathToGenerateTemps []string, contentToGenerateTemps []string) ([]string, []string) {
	validatorHelperPath := strings.TrimSuffix(path, "handler.json") + "validator.go"
	validationHelperContent := handlertp.ValidationHelperTemplate
	validationHelperContent = myfile.RenamePackageGolangFile(validationHelperContent, packageName)

	pathToGenerateTemps = append(pathToGenerateTemps, validatorHelperPath)
	contentToGenerateTemps = append(contentToGenerateTemps, validationHelperContent)

	return pathToGenerateTemps, contentToGenerateTemps
}

func AddContentForTempUserHandlerAndGeneratedHandler(path string, pathToGenerateTemps []string, contentToGenerateTemps []string, userHandlerContent string, generatedHandlerContent string, ht HandlerTemplateData) ([]string, []string) {
	// generate user handler for temp
	currentPath := strings.TrimSuffix(path, "/handler.json")
	userHandlerTempContent := GenerateTempUserHandlerWithSwagGoSyntax(path, userHandlerContent, ht)
	userHandlerTempContent = strings.Replace(userHandlerTempContent, "#moduleName#", fmt.Sprintf("%s/%s/%s/dto", mygo.GetModuleName(), currentPath, GenTempGenerateFolderAndPackageName(path)), -1)

	userHandlerTempPath := strings.TrimSuffix(path, "handler.json") + fmt.Sprintf("%s/", GenTempGenerateFolderAndPackageName(path)) + mystr.ToSnakeCase(ht.Name) + ".go"
	pathToGenerateTemps = append(pathToGenerateTemps, userHandlerTempPath)
	contentToGenerateTemps = append(contentToGenerateTemps, userHandlerTempContent)

	// generate generated handler for temp
	generateHandlerTempPath := strings.TrimSuffix(path, "handler.json") + fmt.Sprintf("%s/", GenTempGenerateFolderAndPackageName(path)) + "dto/" + mystr.ToSnakeCase(ht.Name) + "_gen.go"
	pathToGenerateTemps = append(pathToGenerateTemps, generateHandlerTempPath)
	contentToGenerateTemps = append(contentToGenerateTemps, generatedHandlerContent)

	return pathToGenerateTemps, contentToGenerateTemps
}

func AddContentForUserHandlerAndGeneratedHandler(path string, pathToGenerates []string, contentToGenerates []string, userHandlerContent string, generatedHandlerContent string, ht HandlerTemplateData) ([]string, []string) {
	// generate user handler
	currentPath := strings.TrimSuffix(path, "/handler.json")
	userHandlerContent = strings.Replace(userHandlerContent, "#swaggo#", "", -1)
	userHandlerContent = strings.Replace(userHandlerContent, "#moduleName#", fmt.Sprintf("%s/%s/dto", mygo.GetModuleName(), currentPath), -1)

	userHandlerPath := strings.TrimSuffix(path, "handler.json") + mystr.ToSnakeCase(ht.Name) + ".go"
	pathToGenerates = append(pathToGenerates, userHandlerPath)
	contentToGenerates = append(contentToGenerates, userHandlerContent)

	// generate generated handler
	generateHandlerPath := strings.TrimSuffix(path, "handler.json") + "dto/" + mystr.ToSnakeCase(ht.Name) + "_gen.go"
	pathToGenerates = append(pathToGenerates, generateHandlerPath)
	contentToGenerates = append(contentToGenerates, generatedHandlerContent)

	return pathToGenerates, contentToGenerates
}

func AddContentForMainHandlerAndRouteFile(path string, packageName string, pathToGenerates []string, contentToGenerates []string, htds []HandlerTemplateData) ([]string, []string) {
	// generate main handler file
	mainHandlerContent := GenerateContentMainHandler(packageName)
	mainHandlerOutputPath := strings.TrimSuffix(path, "handler.json") + "handler.go"
	pathToGenerates = append(pathToGenerates, mainHandlerOutputPath)
	contentToGenerates = append(contentToGenerates, mainHandlerContent)

	// generate handler routes file
	routeContents := GenerateContentRoutes(htds, packageName)
	routeOutputPath := strings.TrimSuffix(path, "handler.json") + "routes_gen.go"
	pathToGenerates = append(pathToGenerates, routeOutputPath)
	contentToGenerates = append(contentToGenerates, routeContents)

	return pathToGenerates, contentToGenerates
}

func AddContentForTempMainHandlerAndRouteFile(path string, packageName string, pathToGenerateTemps []string, contentToGenerateTemps []string, htds []HandlerTemplateData) ([]string, []string) {
	// generate main handler for temp generated
	mainHandlerTempPath := strings.TrimSuffix(path, "handler.json") + fmt.Sprintf("%s/", GenTempGenerateFolderAndPackageName(path)) + "handler.go"
	// this will generate swaggo because temp should have swaggo syntax to generate doc
	mainHandlerTempContent := GenerateTempMainHandler(path, packageName)
	mainHandlerTempContent = myfile.RenamePackageGolangFile(mainHandlerTempContent, GenTempGenerateFolderAndPackageName(path))
	pathToGenerateTemps = append(pathToGenerateTemps, mainHandlerTempPath)
	contentToGenerateTemps = append(contentToGenerateTemps, mainHandlerTempContent)

	// generate handler routes file
	routeContents := GenerateContentRoutes(htds, packageName)
	routeOutputPath := strings.TrimSuffix(path, "handler.json") + fmt.Sprintf("%s/", GenTempGenerateFolderAndPackageName(path)) + "routes_gen.go"
	routeContents = myfile.RenamePackageGolangFile(routeContents, GenTempGenerateFolderAndPackageName(path))
	pathToGenerateTemps = append(pathToGenerateTemps, routeOutputPath)
	contentToGenerateTemps = append(contentToGenerateTemps, routeContents)

	return pathToGenerateTemps, contentToGenerateTemps
}

func GenerateContentBothUserHandlerAndGeneratedHandler(packageName string, htd HandlerTemplateData) (string, string) {
	template := handlertp.Template

	template = strings.Replace(template, "handlertp", packageName, -1)
	template = strings.Replace(template, "#handlerName#", mystr.CapitalizeFirstLetter(htd.Name)+"Handler", -1)
	template = strings.Replace(template, "#handlerFuncName#", mystr.CapitalizeFirstLetter(htd.Name), -1)
	template = strings.Replace(template, "#requestName#", htd.Name+"Request", -1)
	template = strings.Replace(template, "#responseName#", htd.Name+"Response", -1)

	templateGenerate := handlertp.TemplateGenerate
	templateGenerate = strings.Replace(templateGenerate, "handlertp", "dto", -1)
	templateGenerate = strings.Replace(templateGenerate, "#handlerFuncName#", mystr.CapitalizeFirstLetter(htd.Name), -1)

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
	template = strings.Replace(template, "handlertp", packageName, -1)
	return template
}

func GenerateContentRoutes(htds []HandlerTemplateData, packageName string) string {
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
	routeTemplate = strings.Replace(routeTemplate, "handlertp", packageName, -1)
	return strings.Replace(routeTemplate, "#route#", routeContents, -1)
}
