package main

import (
	"log"

	"metrics-bench-suite/pkg/cmd/loader"
)

func main() {

	var rootCmd = loader.NewCommand()

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
