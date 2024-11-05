package handlergen

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"
)

var (
	statusCodeStyleTemplateWithRequest = `
	[
	  {
		"type": "handler",
		"name": "createProduct",
		"api": "/products",
		"method": "post",
		"header": "",
		"request": {
		  "type": "string|required",
		  "price": "int"
		},
		"response": {
		  "200": {
			"success": "bool",
			"code": "int",
			"data": {
			  "type": "string",
			  "price": "int",
			  "meet": {
				"type": "string",
				"price": "int"
			  }
			}
		  },
		  "400": {
			"type": "string",
			"many": {
			  "type": "string",
			  "price": "int"
			},
			"price": "int"
		  }
		}
	  }
	]`
)

func TestGenerateHandlerWithStatusCodeStyleTemplateWithRequest(t *testing.T) {
	var handlerTemplates []HandlerTemplateData
	err := json.Unmarshal([]byte(statusCodeStyleTemplateWithRequest), &handlerTemplates)
	if err != nil {
		t.Error("wrong handler template data")
	}
	template, templateGenerate := GenerateHandler("handler", handlerTemplates[0])
	outputTemplatePath := "test/template_generated.go"
	outputTemplateGeneratePath := "test/template_generate_generated.go"

	err = os.WriteFile(outputTemplatePath, []byte(template), 0644)
	if err != nil {
		t.Error("Error writing file", outputTemplatePath, err)
	}

	cmd := exec.Command("go", "fmt", outputTemplatePath)
	err = cmd.Run()
	if err != nil {
		t.Error("Error formatting file", outputTemplatePath, err)
	}

	err = os.WriteFile(outputTemplateGeneratePath, []byte(templateGenerate), 0644)
	if err != nil {
		t.Error("Error writing file", outputTemplateGeneratePath, err)
	}

	cmd = exec.Command("go", "fmt", outputTemplateGeneratePath)
	err = cmd.Run()
	if err != nil {
		t.Error("Error formatting file", outputTemplateGeneratePath, err)
	}

	fileOutputTemplatePath, err := os.ReadFile(outputTemplatePath)
	if err != nil {
		return
	}

	fileOutputTemplateGeneratePath, err := os.ReadFile(outputTemplateGeneratePath)
	if err != nil {
		return
	}

	expectTemplatePath := "test/template.go"
	fileExpectTemplatePath, err := os.ReadFile(expectTemplatePath)
	if err != nil {
		return
	}

	expectTemplateGeneratePath := "test/template_generate.go"
	fileExpectTemplateGeneratePath, err := os.ReadFile(expectTemplateGeneratePath)
	if err != nil {
		return
	}

	if string(fileOutputTemplatePath) != string(fileExpectTemplatePath) {
		t.Error("test error output template mismatch")
	}

	if string(fileOutputTemplateGeneratePath) != string(fileExpectTemplateGeneratePath) {
		t.Error("test error output template generate mismatch")
	}

	err = os.Remove(outputTemplatePath)
	if err != nil {
		t.Error("Error removing file", outputTemplatePath, err)
	}
	err = os.Remove(outputTemplateGeneratePath)
	if err != nil {
		t.Error("Error removing file", outputTemplatePath, err)
	}
}
