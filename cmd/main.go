package main

import (
	"fmt"
	"os"

	"github.com/Petroviiic/GoGit/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gogit <command>")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "init":
		if err := commands.RunInit(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

}
