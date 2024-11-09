package gencmd

import (
	"fmt"
	"github.com/mrbryside/go-generate/internal/generator/handlergen"
	"os"
)

// teardown will write knab_logs.json, remove all swaggo folder, and delete file not in knab_logs.json
func tearDown(dir string, paths []string, handlerGenPathForLogs []string) {
	// delete file not in knab_logs.json
	for _, p := range paths {
		if p == "" {
			continue
		}
		err := deleteFileNotInKnabLogs(p, handlerGenPathForLogs)
		if err != nil {
			fmt.Printf("Error deleting file not in knab_logs.json: %v\n", err)
			return
		}
	}

	// write knab_logs.json for delete unused handler generated
	err := writeKnabLogs(dir, handlerGenPathForLogs)
	if err != nil {
		fmt.Printf("Error writing knab logs: %v\n", err)
		return
	}

	// remove all temp generated folder
	for _, p := range paths {
		if p == "" {
			continue
		}
		err = os.RemoveAll(p + fmt.Sprintf("/%s", handlergen.GenTempGenerateFolderAndPackageName(p)))
		if err != nil {
			fmt.Printf("Error removing folder swaggo: %v\n", err)
			return
		}
	}
}
