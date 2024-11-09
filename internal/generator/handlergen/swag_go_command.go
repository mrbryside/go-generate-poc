package handlergen

import (
	"fmt"
	"os/exec"

	"github.com/mrbryside/go-generate/internal/myfile"
)

func runSwagGoInit(path string) error {
	rootDir := myfile.GetFirstDirectory(path)
	cmd := exec.Command("swag", "init")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v : %s", err, output)
	}
	return nil
}
