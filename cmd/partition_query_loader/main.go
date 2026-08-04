package main

import (
	"log"

	"metrics-bench-suite/pkg/cmd/partition_query_loader"
)

func main() {
	if err := partition_query_loader.NewCommand().Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
