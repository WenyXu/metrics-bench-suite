package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

// Requester is the struct for the requester
type Requester struct {
	URL    string
	Client *http.Client
}

// NewRequester creates a new requester
func NewRequester(url string) *Requester {
	return &Requester{
		URL:    url,
		Client: &http.Client{},
	}
}

// Send sends a write request to the remote write endpoint
func (r *Requester) Send(writeRequest prompb.WriteRequest) error {
	protobufData, err := writeRequest.Marshal()
	if err != nil {
		return err
	}
	compressedData := snappy.Encode(nil, protobufData)
	req, err := http.NewRequest("POST", r.URL, bytes.NewBuffer(compressedData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Content-Encoding", "snappy")

	resp, err := r.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// TODO(weny): check the response
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		return fmt.Errorf("failed to send HTTP request: %v, body: %v", resp.Status, string(body))
	}

	return nil
}
