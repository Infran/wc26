package main

import (
	"os"
)

var Version = "0.1.0"

func main() {
	if err := Execute(); err != nil {
		os.Exit(1)
	}
}
