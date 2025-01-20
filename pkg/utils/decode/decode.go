package decode

import (
	"encoding/base64"

	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

// Base64 decodes a base64 encoded string.
//
// Parameters:
//   - input: The base64 encoded string to decode.
//
// Returns:
//   - []byte: The decoded string.
//   - error: An error if the decoding fails.
func Base64(input string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(input)
}

// Body decodes the body of a Prometheus write request.
//
// Parameters:
//   - compressedData: The compressed data to decode.
//
// Returns:
//   - prompb.WriteRequest: The decoded write request.
//   - error: An error if the decoding fails.
func Body(compressedData []byte) (prompb.WriteRequest, error) {
	protobufData, err := snappy.Decode(nil, compressedData)
	if err != nil {
		return prompb.WriteRequest{}, err
	}

	var writeRequest prompb.WriteRequest
	if err := writeRequest.Unmarshal(protobufData); err != nil {
		return prompb.WriteRequest{}, err
	}

	return writeRequest, nil
}
