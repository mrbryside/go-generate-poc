package handlergen

import (
	"fmt"
	"github.com/mrbryside/go-generate/internal/utils/myfile"
	"os/exec"
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
