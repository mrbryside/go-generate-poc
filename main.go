package main

import (
	"fmt"
	"os"

	"github.com/mrbryside/go-generate/internal/comand/gencmd"
)

func main() {
	if len(os.Args) < 2 {
		println("Expected 'generate' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		gencmd.Run()
		println("[go-generate] Successfully generated files!")
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
