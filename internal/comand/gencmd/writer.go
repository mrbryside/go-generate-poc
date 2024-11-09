package gencmd

import (
	"encoding/json"
	"github.com/mrbryside/go-generate/internal/myfile"
	"os"
)

func writeKnabLogs(path string, allHandlerPathGens []string) error {
	// write knab logs
	knabLogsPath := path + "/knab_logs.json"
	allPathJson, err := json.Marshal(allHandlerPathGens)
	if err != nil {
		return err
	}
	err = os.WriteFile(knabLogsPath, allPathJson, 0644)
	if err != nil {
		return err
	}
	return nil
}

func deleteFileNotInKnabLogs(path string, allHandlerPaths []string) error {
	rootDir := myfile.GetFirstDirectory(path)
	knabLogsPath := rootDir + "/knab_logs.json"
	knabLogData, err := os.ReadFile(knabLogsPath)
	if err != nil {
		return nil
	}

	var allPathKnabs []string
	err = json.Unmarshal(knabLogData, &allPathKnabs)
	if err != nil {
		return nil
	}

	var deletePaths []string
	for _, pathKnab := range allPathKnabs {
		found := false
		for _, handlerPath := range allHandlerPaths {
			if pathKnab == handlerPath {
				found = true
			}
		}
		if !found {
			deletePaths = append(deletePaths, pathKnab)
		}
	}

	err = myfile.DeleteFileByPaths(deletePaths)
	if err != nil {
		return err
	}
	return nil
}
