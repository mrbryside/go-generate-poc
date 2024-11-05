package myfile

import (
	"os"
	"strings"
)

type NewStruct struct {
	Name   string
	Fields string
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

func CreateNewStructs(newStructs []NewStruct) string {
	content := ""
	for _, ns := range newStructs {
		content += "type " + ns.Name + " struct {\n" + ns.Fields + "\n}\n\n"
	}
	return content
}

func AddStructToLastLine(currentContent, fields, structName string) string {
	newContent := CreateNewStructs([]NewStruct{
		{
			Name:   structName,
			Fields: fields,
		},
	})
	return currentContent + "\n" + newContent
}

func AddContentToLastLine(currentContent, newContent string) string {
	return currentContent + "\n" + newContent
}
