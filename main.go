package main

import (
	"fmt"
	"os"
)

func main() {
	if err := RunTUI(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
