package remote_write_request_viewer

import (
	"encoding/base64"
	"errors"
	"log"
	"metrics-bench-suite/pkg/utils/decode"
	"metrics-bench-suite/pkg/utils/format"
	"os"

	"github.com/spf13/cobra"
)

// InputType represents the type of input data
type InputType int

const (
	// Base64Input represents the input type is base64 encoded string
	Base64Input InputType = iota
	// FileInput represents the input type is a file path
	FileInput
)

// Viewer holds the input and output types for the command
type Viewer struct {
	inputType InputType
}

// Run executes the decoding logic
func (dc *Viewer) Run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal("You must provide input data")
	}

	var compressedData []byte
	var err error

	switch dc.inputType {
	case Base64Input:
		inputData := args[0]
		compressedData, err = base64.StdEncoding.DecodeString(inputData)

		if err != nil {
			log.Fatalf("Failed to decode base64 data: %v", err)
		}
	case FileInput:
		inputPath := args[0]
		inputData, err := os.ReadFile(inputPath)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		compressedData, err = base64.StdEncoding.DecodeString(string(inputData))
		if err != nil {
			log.Fatalf("Failed to decode base64 data: %v", err)
		}
	default:
		log.Fatal("Unsupported input type")
	}

	writeRequest, err := decode.Body(compressedData)
	if err != nil {
		log.Fatalf("Failed to decode body: %v", err)
	}

	format.WriteRequest(&writeRequest)
}

func parseInputType(input string) (InputType, error) {
	switch input {
	case "base64":
		return Base64Input, nil
	case "file":
		return FileInput, nil
	default:
		return 0, errors.New("invalid input type")
	}
}

func NewCommand() *cobra.Command {
	viewer := &Viewer{}

	var rootCmd = &cobra.Command{
		Use:   "remote_write_request_viewer <base64_encoded_data> | <file_path>",
		Short: "Remote Write Request Viewer is a tool to decode Prometheus remote write body",
		Run: func(cmd *cobra.Command, args []string) {
			inputTypeStr, _ := cmd.Flags().GetString("type")
			var err error
			viewer.inputType, err = parseInputType(inputTypeStr)
			if err != nil {
				log.Fatal(err)
			}

			viewer.Run(cmd, args)
		},
	}

	rootCmd.Flags().StringP("type", "t", "base64", "Specify the input type (base64 or file)")

	return rootCmd
}
