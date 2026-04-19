package main

import (
	"fmt"
	"os"

	"github.com/Petroviiic/GoGit/internal/commands"
	"github.com/Petroviiic/GoGit/internal/core"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gogit <command>")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	repo, err := core.NewRepository(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "init":
		if err := commands.RunInit(repo); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "add":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Nothing specified, nothing added.\n")
			os.Exit(1)
		}
		if err := commands.RunAdd(args, repo); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

}
