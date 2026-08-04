package partition_query_loader

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"
)

// RegionSnapshot is one cumulative statistics reading for a discovered region.
type RegionSnapshot struct {
	Timestamp          time.Time
	Database           string
	LogicalTable       string
	PhysicalTable      string
	Partition          string
	RegionID           uint64
	LeaderDatanode     string
	DatanodeID         string
	QueryCPUTimeMillis uint64
	QueryScannedBytes  uint64
}

// RequestMetrics contains workload measurements collected during a sampling interval.
type RequestMetrics struct {
	RequestCount uint64
	ErrorCount   uint64
	LatencyP95   time.Duration
}

// Observation records a region's cumulative counters and interval rates.
type Observation struct {
	Timestamp                  time.Time     `json:"timestamp"`
	Database                   string        `json:"database"`
	LogicalTable               string        `json:"logical_table"`
	PhysicalTable              string        `json:"physical_table"`
	Partition                  string        `json:"partition"`
	RegionID                   uint64        `json:"region_id"`
	LeaderDatanode             string        `json:"leader_datanode"`
	DatanodeID                 string        `json:"datanode_id"`
	QueryCPUTimeMillis         uint64        `json:"query_cpu_time_millis"`
	QueryCPUTimeRateCores      float64       `json:"query_cpu_time_rate_cores"`
	QueryScannedBytes          uint64        `json:"query_scanned_bytes"`
	QueryScannedBytesPerSecond float64       `json:"query_scanned_bytes_per_second"`
	RequestCount               uint64        `json:"request_count"`
	ErrorCount                 uint64        `json:"error_count"`
	LatencyP95                 time.Duration `json:"latency_p95"`
	Moved                      bool          `json:"moved"`
	Reset                      bool          `json:"reset"`
	CPUTimeRateValid           bool          `json:"query_cpu_time_rate_valid"`
}

// DatanodeCPUStats summarizes valid region CPU rates for a current datanode topology.
type DatanodeCPUStats struct {
	DatanodeID  string
	RegionCount int
	TotalCores  float64
}

// SampleRegionStatistics reads the cumulative statistics and current leader for every discovered region.
func SampleRegionStatistics(ctx context.Context, db *sql.DB, tables []DiscoveredTable, now time.Time) ([]RegionSnapshot, error) {
	snapshots := make([]RegionSnapshot, 0)
	seen := make(map[uint64]struct{})
	for _, table := range tables {
		for _, partition := range table.Partitions {
			if _, ok := seen[partition.RegionID]; ok {
				continue
			}
			seen[partition.RegionID] = struct{}{}
			rows, err := db.QueryContext(ctx, `SELECT statistics.region_id, peers.peer_id, statistics.query_cpu_time_millis, statistics.query_scanned_bytes FROM information_schema.region_statistics AS statistics LEFT JOIN information_schema.region_peers AS peers ON statistics.region_id = peers.region_id AND peers.is_leader = 'Yes' WHERE statistics.region_id = ?`, partition.RegionID)
			if err != nil {
				return nil, fmt.Errorf("read statistics for region %d: %w", partition.RegionID, err)
			}
			if rows.Next() {
				var snapshot RegionSnapshot
				var leader sql.NullString
				if err := rows.Scan(&snapshot.RegionID, &leader, &snapshot.QueryCPUTimeMillis, &snapshot.QueryScannedBytes); err != nil {
					rows.Close()
					return nil, fmt.Errorf("scan statistics for region %d: %w", partition.RegionID, err)
				}
				if leader.Valid {
					snapshot.LeaderDatanode = leader.String
					snapshot.DatanodeID = leader.String
				}
				snapshot.Timestamp = now
				snapshot.Database = table.Database
				snapshot.LogicalTable = table.LogicalTable
				snapshot.PhysicalTable = table.PhysicalTable
				snapshot.Partition = partition.Name
				snapshots = append(snapshots, snapshot)
			}
			if err := rows.Close(); err != nil {
				return nil, err
			}
			if err := rows.Err(); err != nil {
				return nil, err
			}
		}
	}
	return snapshots, nil
}

// CalculateObservations derives non-negative interval rates and emits moved records for disappeared regions.
func CalculateObservations(current, previous []RegionSnapshot, metrics RequestMetrics) []Observation {
	prior := make(map[uint64]RegionSnapshot, len(previous))
	for _, snapshot := range previous {
		prior[snapshot.RegionID] = snapshot
	}
	seen := make(map[uint64]struct{}, len(current))
	observations := make([]Observation, 0, len(current)+len(previous))
	for _, snapshot := range current {
		seen[snapshot.RegionID] = struct{}{}
		observation := observationFromSnapshot(snapshot, metrics)
		if old, ok := prior[snapshot.RegionID]; ok {
			if snapshotDatanode(old) != snapshotDatanode(snapshot) {
				observation.Moved = true
			} else if snapshot.QueryCPUTimeMillis < old.QueryCPUTimeMillis || snapshot.QueryScannedBytes < old.QueryScannedBytes {
				observation.Reset = true
			} else if elapsed := snapshot.Timestamp.Sub(old.Timestamp); elapsed > 0 {
				seconds := elapsed.Seconds()
				observation.QueryCPUTimeRateCores = float64(snapshot.QueryCPUTimeMillis-old.QueryCPUTimeMillis) / (elapsed.Seconds() * 1000)
				observation.QueryScannedBytesPerSecond = float64(snapshot.QueryScannedBytes-old.QueryScannedBytes) / seconds
				observation.CPUTimeRateValid = true
			}
		}
		observations = append(observations, observation)
	}
	for regionID, snapshot := range prior {
		if _, ok := seen[regionID]; !ok {
			observation := observationFromSnapshot(snapshot, metrics)
			observation.Moved = true
			observations = append(observations, observation)
		}
	}
	return observations
}

func snapshotDatanode(snapshot RegionSnapshot) string {
	if snapshot.DatanodeID != "" {
		return snapshot.DatanodeID
	}
	return snapshot.LeaderDatanode
}

func observationFromSnapshot(snapshot RegionSnapshot, metrics RequestMetrics) Observation {
	datanode := snapshot.DatanodeID
	if datanode == "" {
		datanode = snapshot.LeaderDatanode
	}
	return Observation{Timestamp: snapshot.Timestamp, Database: snapshot.Database, LogicalTable: snapshot.LogicalTable, PhysicalTable: snapshot.PhysicalTable, Partition: snapshot.Partition, RegionID: snapshot.RegionID, LeaderDatanode: datanode, DatanodeID: datanode, QueryCPUTimeMillis: snapshot.QueryCPUTimeMillis, QueryScannedBytes: snapshot.QueryScannedBytes, RequestCount: metrics.RequestCount, ErrorCount: metrics.ErrorCount, LatencyP95: metrics.LatencyP95}
}

// AggregateDatanodeCPU groups current regions and valid CPU rates by current peer ID.
func AggregateDatanodeCPU(current, previous []RegionSnapshot) []DatanodeCPUStats {
	observations := CalculateObservations(current, previous, RequestMetrics{})
	stats := make(map[string]*DatanodeCPUStats)
	for _, observation := range observations {
		datanode := observation.LeaderDatanode
		if datanode == "" {
			continue
		}
		item := stats[datanode]
		if item == nil {
			item = &DatanodeCPUStats{DatanodeID: datanode}
			stats[datanode] = item
		}
		item.RegionCount++
		if observation.CPUTimeRateValid {
			item.TotalCores += observation.QueryCPUTimeRateCores
		}
	}
	result := make([]DatanodeCPUStats, 0, len(stats))
	for _, item := range stats {
		result = append(result, *item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].DatanodeID < result[j].DatanodeID })
	return result
}
