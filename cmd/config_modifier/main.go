package main

import (
	"log"

	"metrics-bench-suite/pkg/cmd/config_modifier"
)

func main() {

	var rootCmd = config_modifier.NewCommand()

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
