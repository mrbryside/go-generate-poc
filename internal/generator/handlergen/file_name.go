package handlergen

import "strings"

func GenBasePath(path string) string {
	return strings.TrimSuffix(path, "/handler.json")
}
