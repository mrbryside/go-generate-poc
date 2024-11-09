package gencmd

import (
	"encoding/json"
	"github.com/mrbryside/go-generate/internal/generator/handlergen"
	"os"
)

func validateUnmarshalHandlerJSONSpec(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var handlerTemplates []handlergen.HandlerTemplateData
	err = json.Unmarshal(data, &handlerTemplates)
	if err != nil {
		return err
	}
	return nil
}
