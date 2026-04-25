package main

import (
	"fmt"
	"os"

	"github.com/Petroviiic/GoGit/internal/commands"
	"github.com/Petroviiic/GoGit/internal/core"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./... <command>")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	repo, err := core.NewRepository(args)

	if err != nil && command != "init" {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "test":
		//add anything
		commands.TestFunc(repo)

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
	case "commit":
		var message, author string
		author = "GoGit User <user@gogit.com>"

		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "-m":
				if i+1 < len(args) {
					if args[i+1] == "--author" {
						fmt.Fprintln(os.Stderr, "Error: Parameters malformed")
						os.Exit(1)
					}
					message = args[i+1]
					i++
				}
			case "--author":
				if i+1 < len(args) {
					if args[i+1] == "-m" {
						fmt.Fprintln(os.Stderr, "Error: Parameters malformed")
						os.Exit(1)
					}
					author = args[i+1]
					i++
				}
			}
		}

		if message == "" {
			fmt.Fprintln(os.Stderr, "Error: Commit message is required (-m)")
			os.Exit(1)
		}

		if err := commands.RunCommit(repo, message, author); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "checkout":
		if len(args) > 2 {
			fmt.Fprintln(os.Stderr, "Error: Parameters malformed")
			os.Exit(1)
		}

		shouldCreate := false
		if len(args) == 2 {
			if args[0] != "-b" {
				fmt.Fprintln(os.Stderr, "Error: Parameters malformed")
				os.Exit(1)
			} else {
				shouldCreate = true
			}
		}

		branch := args[len(args)-1]

		if err := commands.RunCheckout(branch, shouldCreate, repo); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

}
