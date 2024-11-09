package handlergen

import (
	"os/exec"

	"github.com/mrbryside/go-generate/internal/myfile"
)

func runSwagGoInit(path string) error {
	rootDir := myfile.GetFirstDirectory(path)
	cmd := exec.Command("swag", "init")
	cmd.Dir = rootDir
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
