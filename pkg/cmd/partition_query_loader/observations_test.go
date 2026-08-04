package partition_query_loader

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"strings"
	"testing"
	"time"
)

func TestCalculateObservationsRates(t *testing.T) {
	start := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	observations := CalculateObservations([]RegionSnapshot{snapshot(start.Add(2*time.Second), 1, "node-1", 250, 900)}, []RegionSnapshot{snapshot(start, 1, "node-1", 50, 100)}, RequestMetrics{RequestCount: 4, ErrorCount: 1, LatencyP95: 12 * time.Millisecond})
	if len(observations) != 1 {
		t.Fatalf("observations = %d", len(observations))
	}
	got := observations[0]
	if got.QueryCPUTimeRateCores != 0.1 || got.QueryScannedBytesPerSecond != 400 || got.RequestCount != 4 || got.ErrorCount != 1 || got.LatencyP95 != 12*time.Millisecond {
		t.Fatalf("unexpected observation: %#v", got)
	}
}

func TestCalculateObservationsFirstAndZeroElapsedSamples(t *testing.T) {
	now := time.Now()
	first := CalculateObservations([]RegionSnapshot{snapshot(now, 1, "node-1", 20, 200)}, nil, RequestMetrics{})[0]
	zeroElapsed := CalculateObservations([]RegionSnapshot{snapshot(now, 1, "node-1", 30, 300)}, []RegionSnapshot{snapshot(now, 1, "node-1", 20, 200)}, RequestMetrics{})[0]
	if first.QueryCPUTimeRateCores != 0 || first.QueryScannedBytesPerSecond != 0 || zeroElapsed.QueryCPUTimeRateCores != 0 || zeroElapsed.QueryScannedBytesPerSecond != 0 {
		t.Fatalf("first=%#v zeroElapsed=%#v", first, zeroElapsed)
	}
}

func TestCalculateObservationsSupportsSubMillisecondIntervals(t *testing.T) {
	start := time.Now()
	got := CalculateObservations([]RegionSnapshot{snapshot(start.Add(500*time.Microsecond), 1, "node-1", 1, 1)}, []RegionSnapshot{snapshot(start, 1, "node-1", 0, 0)}, RequestMetrics{})[0]
	if got.QueryCPUTimeRateCores != 2 || got.QueryScannedBytesPerSecond != 2000 {
		t.Fatalf("unexpected sub-millisecond rates: %#v", got)
	}
}

func TestCalculateObservationsDeduplicatesPhysicalRegions(t *testing.T) {
	now := time.Now()
	duplicates := []RegionSnapshot{snapshot(now, 1, "node-1", 1, 1), snapshot(now, 1, "node-1", 1, 1)}
	if observations := observationsFromSnapshots([][]RegionSnapshot{duplicates}, nil); len(observations) != 1 {
		t.Fatalf("observations=%#v", observations)
	}
}

func TestCalculateObservationsMarksResetMovementAndDisappearance(t *testing.T) {
	now := time.Now()
	previous := []RegionSnapshot{snapshot(now, 1, "node-1", 100, 100), snapshot(now, 2, "node-1", 100, 100), snapshot(now, 3, "node-1", 100, 100)}
	current := []RegionSnapshot{snapshot(now.Add(time.Second), 1, "node-1", 10, 200), snapshot(now.Add(time.Second), 2, "node-2", 200, 200)}
	observations := CalculateObservations(current, previous, RequestMetrics{})
	byRegion := make(map[uint64]Observation)
	for _, observation := range observations {
		byRegion[observation.RegionID] = observation
	}
	if !byRegion[1].Reset || byRegion[1].QueryCPUTimeRateCores != 0 || byRegion[1].QueryScannedBytesPerSecond != 0 {
		t.Fatalf("reset not handled: %#v", byRegion[1])
	}
	if !byRegion[2].Moved || byRegion[2].QueryCPUTimeRateCores != 0 || byRegion[2].QueryScannedBytesPerSecond != 0 {
		t.Fatalf("movement not handled: %#v", byRegion[2])
	}
	if !byRegion[3].Moved || byRegion[3].Reset {
		t.Fatalf("disappearance not handled: %#v", byRegion[3])
	}
}

func TestAggregateDatanodeCPURefreshesTopologyAndExcludesMovedDelta(t *testing.T) {
	start := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	previous := []RegionSnapshot{snapshot(start, 1, "node-a", 100, 100), snapshot(start, 2, "node-a", 100, 100)}
	current := []RegionSnapshot{snapshot(start.Add(time.Second), 1, "node-b", 300, 300), snapshot(start.Add(time.Second), 2, "node-a", 300, 300)}
	stats := AggregateDatanodeCPU(current, previous)
	if len(stats) != 2 || stats[0].DatanodeID != "node-a" || stats[0].RegionCount != 1 || stats[0].TotalCores != 0.2 || stats[1].DatanodeID != "node-b" || stats[1].RegionCount != 1 || stats[1].TotalCores != 0 {
		t.Fatalf("unexpected datanode stats: %#v", stats)
	}
}

func TestSampleRegionStatistics(t *testing.T) {
	db := mockStatisticsDB(t)
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	tables := []DiscoveredTable{{Database: "metrics", LogicalTable: "cpu", PhysicalTable: "physical", Partitions: []DiscoveredPartition{{Name: "p0", RegionID: 42}}}}
	snapshots, err := SampleRegionStatistics(context.Background(), db, tables, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshots) != 1 || snapshots[0].RegionID != 42 || snapshots[0].LeaderDatanode != "node-1" || snapshots[0].DatanodeID != "node-1" || snapshots[0].QueryCPUTimeMillis != 12 || snapshots[0].QueryScannedBytes != 99 || !snapshots[0].Timestamp.Equal(now) {
		t.Fatalf("unexpected snapshots: %#v", snapshots)
	}
}

func snapshot(timestamp time.Time, regionID uint64, leader string, cpu, bytes uint64) RegionSnapshot {
	return RegionSnapshot{Timestamp: timestamp, Database: "metrics", LogicalTable: "cpu", PhysicalTable: "physical", Partition: "p0", RegionID: regionID, LeaderDatanode: leader, QueryCPUTimeMillis: cpu, QueryScannedBytes: bytes}
}

type statisticsDriver struct{}
type statisticsConn struct{}
type statisticsRows struct{ returned bool }

func (statisticsDriver) Open(string) (driver.Conn, error)  { return statisticsConn{}, nil }
func (statisticsConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (statisticsConn) Close() error                        { return nil }
func (statisticsConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }
func (statisticsConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	if !strings.Contains(query, "information_schema.region_statistics") || !strings.Contains(query, "region_peers") || !strings.Contains(query, "is_leader = 'Yes'") {
		return nil, driver.ErrSkip
	}
	return &statisticsRows{}, nil
}
func (r *statisticsRows) Columns() []string {
	return []string{"region_id", "peer_id", "query_cpu_time_millis", "query_scanned_bytes"}
}
func (r *statisticsRows) Close() error { return nil }
func (r *statisticsRows) Next(dest []driver.Value) error {
	if r.returned {
		return io.EOF
	}
	r.returned = true
	dest[0], dest[1], dest[2], dest[3] = int64(42), "node-1", int64(12), int64(99)
	return nil
}
func mockStatisticsDB(t *testing.T) *sql.DB {
	t.Helper()
	name := "statistics_mock_" + strings.ReplaceAll(t.Name(), "/", "_")
	sql.Register(name, statisticsDriver{})
	db, err := sql.Open(name, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
