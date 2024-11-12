package handlergen

import (
	"fmt"
	"github.com/mrbryside/go-generate/internal/utils/mystr"
	"strings"
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

func GenerateGoFileNameInBasePath(path string, fileName string) string {
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
	newPath := fmt.Sprintf("%s/%s", basePath, DtoFolderAndPackageName)
	return GenerateGoFileNameInPath(newPath, newFileName)
}

func GenHandlerFileNameFromHandlerTemplate(htd HandlerTemplateData) string {
	return mystr.ToSnakeCase(htd.Name)
}
