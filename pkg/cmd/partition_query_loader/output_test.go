package partition_query_loader

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWriteObservationsNDJSON(t *testing.T) {
	var output bytes.Buffer
	observation := observationFromSnapshot(snapshot(time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC), 42, "node-1", 10, 20), RequestMetrics{RequestCount: 3})
	if err := WriteObservations(&output, "ndjson", []Observation{observation}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "\"region_id\":42") || !strings.Contains(output.String(), "\"request_count\":3") || !strings.HasSuffix(output.String(), "\n") {
		t.Fatalf("unexpected NDJSON: %q", output.String())
	}
}

func TestWriteObservationsCSV(t *testing.T) {
	var output bytes.Buffer
	observation := observationFromSnapshot(snapshot(time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC), 42, "node-1", 10, 20), RequestMetrics{})
	if err := WriteObservations(&output, "csv", []Observation{observation}); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) != 2 || !strings.Contains(lines[0], "query_cpu_time_rate_cores") || !strings.Contains(lines[1], ",42,node-1,node-1,10,") {
		t.Fatalf("unexpected CSV: %q", output.String())
	}
}
