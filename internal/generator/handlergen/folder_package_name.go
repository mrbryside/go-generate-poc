package handlergen

import (
	"strings"

	"github.com/mrbryside/go-generate/internal/generator/template/handlertp"
)

func GenTempGenerateFolderAndPackageName(path string) string {
	result := GenBasePath(path)
	result = strings.Replace(result, "/", "_", -1)
	result += "_" + handlertp.TempGenerateFolderAndPackageName
	return result
}
