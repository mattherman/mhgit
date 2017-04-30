package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Too few arguments provided.")
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		fmt.Println("git init")
	}
}
