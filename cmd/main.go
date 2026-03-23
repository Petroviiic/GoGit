package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		if args[1] == "init" {
			fmt.Println("gitinit")
		}
	}

}
