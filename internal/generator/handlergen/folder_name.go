package handlergen

import (
	"strings"
)

const (
	// TempGenerateFolderAndPackageName value `entities` will show in swagger object
	TempGenerateFolderAndPackageName = "entities"
)

func GenTempGenerateFolderAndPackageName(path string) string {
	result := GenBasePath(path)
	result = strings.Replace(result, "/", "_", -1)
	return result
}
