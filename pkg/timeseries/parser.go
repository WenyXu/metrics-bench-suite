package timeseries

import (
	"log"
	"metrics-bench-suite/pkg/parser"
	"metrics-bench-suite/pkg/utils/decode"

	"github.com/prometheus/prometheus/prompb"
)

// GetUniqueTimeSeries is the function to get the unique time series from remote write requests
func GetUniqueTimeSeries(wrSet []prompb.WriteRequest) ([]TimeSerie, error) {
	counter := make(Counter)
	for _, wr := range wrSet {
		err := CountAllTimeSeries(&wr, &counter)
		if err != nil {
			return nil, err
		}
	}

	tsSet, err := counter.GetAllTimeSeries()
	log.Printf("Found %d time series", len(tsSet))
	if err != nil {
		return nil, err
	}

	return tsSet, nil
}

// ParseAllRemoteWriteRequest is the function to parse the all remote write request from the tcpflow output
func ParseAllRemoteWriteRequest(tcpflowOutput string) ([]prompb.WriteRequest, error) {
	log.Printf("Parsing request file %s", tcpflowOutput)
	requests, err := parser.ParseHTTPRequests(tcpflowOutput)
	if err != nil {
		log.Printf("Error parsing request file %s: %v", tcpflowOutput, err)
		return nil, err
	}

	wrSet := []prompb.WriteRequest{}
	for _, request := range requests {
		wr, err := decode.Body(request.Body)
		if err != nil {
			log.Fatalf("failed to decode body: %v", err)
		}
		wrSet = append(wrSet, wr)
	}
	return wrSet, nil
}
