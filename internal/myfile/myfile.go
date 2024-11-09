package myfile

import (
	"fmt"
	"github.com/mrbryside/go-generate/internal/mystr"
	"os"
	"path/filepath"
	"strings"
)

type NewStruct struct {
	Name   string
	Fields string
}

func CreateFolderIfNotExist(path string, folderName string) error {
	dirPath := filepath.Join(filepath.Dir(path), folderName)

	// Check if directory already exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Create directory if it doesn't exist
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %w", dirPath, err)
		}
	}
	return nil
}

func GetFirstDirectory(path string) string {
	parts := strings.Split(path, string(os.PathSeparator))
	if len(parts) > 1 {
		return parts[0]
	}
	return ""
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func RemoveLine(template, placeholder string) string {
	lines := strings.Split(template, "\n")
	var result []string

	for _, line := range lines {
		if !strings.Contains(line, placeholder) {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func CreateNewWithSwagGoNameCommentStructs(newStructs []NewStruct) string {
	content := ""
	for _, ns := range newStructs {
		content += "type " + ns.Name + " struct {\n" + ns.Fields + "\n}// @name " + ns.Name + "\n\n"

	}
	return content
}

func AddStructToLastLine(currentContent, fields, structName string) string {
	newContent := CreateNewWithSwagGoNameCommentStructs([]NewStruct{
		{
			Name:   mystr.CapitalizeFirstLetter(structName),
			Fields: fields,
		},
	})
	return currentContent + "\n" + newContent
}

func AddContentToLastLine(currentContent, newContent string) string {
	return currentContent + "\n" + newContent
}

func RenamePackageGolangFile(currentContent, newPackageName string) string {
	lines := strings.Split(currentContent, "\n")
	var result []string

	for _, line := range lines {
		if strings.Contains(line, "package") {
			result = append(result, "package "+newPackageName)
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

func DeleteFileByPaths(paths []string) error {
	for _, p := range paths {
		err := os.RemoveAll(p)
		if err != nil {
			return err
		}
	}
	return nil
}
