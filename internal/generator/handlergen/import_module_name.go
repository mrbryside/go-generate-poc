package handlergen

import (
	"fmt"
	"github.com/mrbryside/go-generate/internal/generator/template/handlertp"
	"github.com/mrbryside/go-generate/internal/utils/mygo"
)

func GenModuleDtoTempPath(path string) string {
	basePath := GenBasePath(path)
	return fmt.Sprintf("%s/%s/%s/%s", mygo.GetModuleName(), basePath, GenTempGenerateFolderAndPackageName(path), handlertp.DtoFolderAndPackageName)
}

func GenModuleDtoPath(path string) string {
	basePath := GenBasePath(path)
	return fmt.Sprintf("%s/%s/%s", mygo.GetModuleName(), basePath, handlertp.DtoFolderAndPackageName)
}
