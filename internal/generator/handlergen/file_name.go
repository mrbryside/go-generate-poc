package handlergen

import (
	"fmt"
	"github.com/mrbryside/go-generate/internal/generator/template/handlertp"
	"strings"
)

const (
	ValidatorFileName = "validator"
	HandlerFileName   = "handler"
	RouteFileName     = "routes_gen"
)

func GenBasePath(path string) string {
	return strings.TrimSuffix(path, "/handler.json")
}

func GenerateGoFileName(fileName string) string {
	return fmt.Sprintf("%s.go", fileName)
}

func GenerateGoFileNameInPath(path string, fileName string) string {
	return GenerateGoFileName(fmt.Sprintf("%s/%s", path, fileName))
}

func GenGoFileNameInBasePath(path string, fileName string) string {
	basePath := GenBasePath(path)
	return GenerateGoFileName(fmt.Sprintf("%s/%s", basePath, fileName))
}

func GenGoFileNameGeneratedInBasePath(path string, fileName string) string {
	basePath := GenBasePath(path)
	newFileName := fmt.Sprintf("%s_gen", fileName)
	return GenerateGoFileNameInPath(basePath, newFileName)
}

func GenGoFileNameGeneratedInDtoBasePath(path string, fileName string) string {
	basePath := GenBasePath(path)
	newFileName := fmt.Sprintf("%s_gen", fileName)
	newPath := fmt.Sprintf("%s/%s", basePath, handlertp.DtoFolderAndPackageName)
	return GenerateGoFileNameInPath(newPath, newFileName)
}

func GenGoFileNameGeneratedInDtoTempBasePath(path string, fileName string) string {
	basePath := GenBasePath(path)
	tempBasePath := fmt.Sprintf("%s/%s", basePath, GenTempGenerateFolderAndPackageName(basePath))
	newFileName := fmt.Sprintf("%s_gen", fileName)
	newPath := fmt.Sprintf("%s/%s", tempBasePath, handlertp.DtoFolderAndPackageName)
	return GenerateGoFileNameInPath(newPath, newFileName)
}

func GenGoFileNameInTempBasePath(path string, fileName string) string {
	basePath := GenBasePath(path)
	tempBasePath := fmt.Sprintf("%s/%s", basePath, GenTempGenerateFolderAndPackageName(basePath))
	return GenGoFileNameInBasePath(tempBasePath, fileName)
}

func GenGoFileNameGeneratedInTempBasePath(path string, fileName string) string {
	basePath := GenBasePath(path)
	tempBasePath := fmt.Sprintf("%s/%s", basePath, GenTempGenerateFolderAndPackageName(basePath))
	newFileName := fmt.Sprintf("%s_gen", fileName)
	return GenGoFileNameInBasePath(tempBasePath, newFileName)
}
