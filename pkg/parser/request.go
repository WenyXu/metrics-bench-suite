package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	// ContentLength is the header key for the content length
	ContentLength = "Content-Length"
)

// HTTPRequest represents a Prometheus remote write request
//
// POST /v1/prometheus/write?db=metrics9 HTTP/1.1
// Host: internal-greptimedb-standalone.internal-greptimedb-standalone.svc.cluster.local:4000
// User-Agent: Prometheus/2.47.2
// Content-Length: 72206
// Content-Encoding: snappy
// Content-Type: application/x-protobuf
// X-Prometheus-Remote-Write-Version: 0.1.0
type HTTPRequest struct {
	Method        string
	URL           string
	Version       string
	ContentLength int
	Headers       map[string]string
	Body          []byte
}

func parseHTTPRequest(scanner *bufio.Reader) (request HTTPRequest, err error) {
	request = HTTPRequest{
		Headers: make(map[string]string),
	}
	requestLine, _, err := scanner.ReadLine()
	if err != nil {
		return request, err
	}
	parts := strings.Fields(string(requestLine))
	if len(parts) < 3 {
		return request, fmt.Errorf("invalid request line")
	}
	request.Method = parts[0]
	request.URL = parts[1]
	request.Version = parts[2]

	for {
		line, _, err := scanner.ReadLine()
		if err != nil {
			return request, err
		}
		if len(line) == 0 {
			break
		}
		parts := strings.SplitN(string(line), ":", 2)
		if len(parts) == 2 {
			request.Headers[parts[0]] = strings.TrimSpace(parts[1])
			if parts[0] == ContentLength {
				request.ContentLength, err = strconv.Atoi(strings.TrimSpace(parts[1]))
				if err != nil {
					return request, fmt.Errorf("invalid content length: %v", parts[1])
				}
			}
		}
	}

	body := make([]byte, request.ContentLength)
	_, err = io.ReadFull(scanner, body)
	if err != nil {
		return request, err
	}
	request.Body = body
	return request, nil
}

// ParseHTTPRequests parses the HTTP requests from the given file.
func ParseHTTPRequests(filePath string) ([]HTTPRequest, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	requests := make([]HTTPRequest, 0)
	for {
		request, err := parseHTTPRequest(reader)
		if err != nil {
			// Ignore EOF and ErrUnexpectedEOF
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}
