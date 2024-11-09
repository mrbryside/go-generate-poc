package gencmd

import (
	"fmt"
	"github.com/mrbryside/go-generate/internal/myfile"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrbryside/go-generate/internal/generator/handlergen"
)

func Run() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Directory argument is required (should be root directory of your api server)")
		os.Exit(1)
	}

	dir := os.Args[2]
	var paths []string
	var handlerGenPathForLogs []string
	var jsonSpecFailed []string
	var reports []handlergen.Report
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			println("Error walking directories:", err.Error())
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "handler.json") {
			err := validateUnmarshalHandlerJSONSpec(path)
			if err != nil {
				println("Error validate handler.json spec:", err.Error())
				return err
			}
			// create temp generated folder
			err = myfile.CreateFolderIfNotExist(path, handlergen.GenTempGenerateFolderAndPackageName(path))
			if err != nil {
				fmt.Printf("Error creating swaggo directory: %v\n", err)
				jsonSpecFailed = append(jsonSpecFailed, path)
				return nil
			}

			// main generate handler
			report := handlergen.MainGenerateHandler(path)
			reports = append(reports, report)
			paths = append(paths, report.BasePathOfJsonSpec)
			handlerGenPathForLogs = append(handlerGenPathForLogs, report.PathToGenerateSuccess...)
		}
		return nil
	})
	if err != nil {
		return
	}

	printReports(reports)
	tearDown(dir, paths, handlerGenPathForLogs)
}
