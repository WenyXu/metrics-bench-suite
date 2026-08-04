package partition_query_loader

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// WriteObservations encodes observations as newline-delimited JSON or CSV.
func WriteObservations(writer io.Writer, format string, observations []Observation) error {
	switch format {
	case "ndjson":
		encoder := json.NewEncoder(writer)
		for _, observation := range observations {
			if err := encoder.Encode(observation); err != nil {
				return err
			}
		}
		return nil
	case "csv":
		return writeCSVObservations(writer, observations)
	default:
		return fmt.Errorf("unsupported observation format %q", format)
	}
}

func writeCSVObservations(writer io.Writer, observations []Observation) error {
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write([]string{"timestamp", "database", "logical_table", "physical_table", "partition", "region_id", "leader_datanode", "datanode_id", "query_cpu_time_millis", "query_cpu_time_rate_cores", "query_cpu_time_rate_valid", "query_scanned_bytes", "query_scanned_bytes_per_second", "request_count", "error_count", "latency_p95", "moved", "reset"}); err != nil {
		return err
	}
	for _, observation := range observations {
		if err := csvWriter.Write([]string{observation.Timestamp.Format("2006-01-02T15:04:05.999999999Z07:00"), observation.Database, observation.LogicalTable, observation.PhysicalTable, observation.Partition, strconv.FormatUint(observation.RegionID, 10), observation.LeaderDatanode, observation.DatanodeID, strconv.FormatUint(observation.QueryCPUTimeMillis, 10), strconv.FormatFloat(observation.QueryCPUTimeRateCores, 'f', -1, 64), strconv.FormatBool(observation.CPUTimeRateValid), strconv.FormatUint(observation.QueryScannedBytes, 10), strconv.FormatFloat(observation.QueryScannedBytesPerSecond, 'f', -1, 64), strconv.FormatUint(observation.RequestCount, 10), strconv.FormatUint(observation.ErrorCount, 10), observation.LatencyP95.String(), strconv.FormatBool(observation.Moved), strconv.FormatBool(observation.Reset)}); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return csvWriter.Error()
}
