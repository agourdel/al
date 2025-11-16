package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alex/al/cmd"
)

func main() {
	// Get the binary name to determine which command to run
	binaryName := filepath.Base(os.Args[0])
	
	// Check if it's a compound command (algo, alnote, allink, alinit)
	if strings.HasPrefix(binaryName, "al") && binaryName != "al" {
		subCommand := strings.TrimPrefix(binaryName, "al")
		
		// Reconstruct args to simulate "al <subcommand> <args>"
		newArgs := []string{os.Args[0], subCommand}
		if len(os.Args) > 1 {
			newArgs = append(newArgs, os.Args[1:]...)
		}
		os.Args = newArgs
	}
	
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
