package format

import (
	"fmt"

	"github.com/prometheus/prometheus/prompb"
)

// WriteRequest formats a Prometheus write request into a human-readable format.
func WriteRequest(writeRequest *prompb.WriteRequest) {
	fmt.Println("Decoded WriteRequest:")
	for i, ts := range writeRequest.Timeseries {
		fmt.Printf("\nTime Series %d:\n", i+1)
		fmt.Println("Labels:")
		for _, label := range ts.Labels {
			fmt.Printf("  %s = %s\n", label.Name, label.Value)
		}
		fmt.Println("Samples:")
		for _, sample := range ts.Samples {
			fmt.Printf("  Value: %f, Timestamp: %d\n", sample.Value, sample.Timestamp)
		}
	}
}
