package main

import (
	"log"
	"metrics-bench-suite/pkg/cmd/schema_generator"
)

func main() {
	var rootCmd = schema_generator.NewCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
