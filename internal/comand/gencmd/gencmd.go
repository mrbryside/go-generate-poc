package gencmd

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mrbryside/go-generate/internal/generator/handlergen"
	"github.com/mrbryside/go-generate/internal/mystr"
)

func Run() {
	dir := "."
	if len(os.Args) > 2 {
		dir = os.Args[2]
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			jsonProcess(path)
		}
		return nil
	})

	if err != nil {
		println("Error walking directories:", err)
	}
}

func jsonProcess(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file %s: %v\n", path, err)
		return
	}

	fileName := strings.TrimSuffix(filepath.Base(path), ".json")
	var contents []string
	var outputPaths []string
	if fileName == "handler" {
		var handlerTemplates []handlergen.HandlerTemplateData
		err = json.Unmarshal(data, &handlerTemplates)
		if err != nil {
			log.Printf("Error format for your handler json: %v\n", err)
			return
		}
		// generate handler routes file
		routeContents := handlergen.GenerateHandlerRoutes(handlerTemplates)
		routeOutputPath := strings.TrimSuffix(path, "handler.json") + "route_gen.go"
		outputPaths = append(outputPaths, routeOutputPath)
		contents = append(contents, routeContents)

		// generate main handler file
		mainHandlerContent := handlergen.GenerateMainHandler()
		mainHandlerOutputPath := strings.TrimSuffix(path, "handler.json") + "handler.go"
		outputPaths = append(outputPaths, mainHandlerOutputPath)
		contents = append(contents, mainHandlerContent)

		// generate all handler file
		for _, ht := range handlerTemplates {
			packageName := filepath.Base(filepath.Dir(path))
			err := handlergen.ValidateHandler(ht)
			if err != nil {
				log.Printf("Error validating your handler json: %v\n", err)
				return
			}
			handlerContent, handlerContentGen := handlergen.GenerateHandler(packageName, ht)
			userHandlerPath := strings.TrimSuffix(path, "handler.json") + mystr.ToSnakeCase(ht.Name) + ".go"
			generateHandlerPath := strings.TrimSuffix(path, "handler.json") + mystr.ToSnakeCase(ht.Name) + "_gen.go"
			//if !myfile.IsFileExist(userHandlerPath) {
			outputPaths = append(outputPaths, userHandlerPath)
			contents = append(contents, handlerContent)
			//}
			outputPaths = append(outputPaths, generateHandlerPath)
			contents = append(contents, handlerContentGen)
		}
	}

	for i, op := range outputPaths {
		err = os.WriteFile(op, []byte(contents[i]), 0644)
		if err != nil {
			log.Printf("Error writing file %s: %v\n", op, err)
			return
		}

		cmd := exec.Command("go", "fmt", op)
		err = cmd.Run()
		if err != nil {
			log.Printf("Error formatting file %s: %v\n", op, err)
		}
	}
}
