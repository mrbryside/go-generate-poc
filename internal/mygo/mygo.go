package mygo

import (
	"os/exec"
	"strings"
)

func GetModuleName() string {
	// Execute the `go list -m` command
	out, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		panic("Error getting your module name: " + err.Error())
	}
	// Convert output to string and trim spaces
	moduleName := strings.TrimSpace(string(out))
	return moduleName
}
