package handlergen

import (
	"strings"
)

const (
	// TempGenerateFolderAndPackageName is the name of the folder and package that will be generated
	TempGenerateFolderAndPackageName = "temp_generated"
	DtoFolderAndPackageName          = "dto"
)

func GenTempGenerateFolderAndPackageName(path string) string {
	result := GenBasePath(path)
	result = strings.Replace(result, "/", "_", -1)
	result += "_" + TempGenerateFolderAndPackageName
	return result
}
