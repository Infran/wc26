package main

import (
	"os"
)

var Version = "0.3.0"

func main() {
	if err := Execute(); err != nil {
		os.Exit(1)
	}
}
