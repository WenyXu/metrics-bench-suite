package partition_query_loader

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
)

func TestSelectTablesIsBoundedAndSeeded(t *testing.T) {
	tables := []string{"a", "b", "c", "d"}
	seed := int64(99)
	left := selectTables(tables, 2, randForTest(seed))
	right := selectTables(tables, 2, randForTest(seed))
	if len(left) != 2 || strings.Join(left, ",") != strings.Join(right, ",") {
		t.Fatalf("selection is not bounded/reproducible: %v vs %v", left, right)
	}
	if got := selectTables(tables, 0, randForTest(seed)); len(got) != len(tables) {
		t.Fatalf("zero limit selected %v", got)
	}
}

func randForTest(seed int64) *rand.Rand { return rand.New(rand.NewSource(seed)) }

func TestDiscoverMapsSafeLogicalTable(t *testing.T) {
	db := mockDiscoveryDB(t, map[string]mockRows{
		"FROM information_schema.tables":         {columns: []string{"table_name"}, values: [][]driver.Value{{"cpu"}}},
		"SHOW CREATE TABLE `metrics`.`cpu`":      {columns: []string{"Table", "Create Table"}, values: [][]driver.Value{{"cpu", "CREATE TABLE cpu ENGINE=metric WITH (on_physical_table = 'physical') PARTITION ON COLUMNS (region) (region = 'east')"}}},
		"SHOW CREATE TABLE `metrics`.`physical`": {columns: []string{"Table", "Create Table"}, values: [][]driver.Value{{"physical", "CREATE TABLE physical PARTITION ON COLUMNS (region) (region = 'east')"}}},
		"FROM information_schema.columns":        {columns: []string{"column_name", "data_type", "column_key"}, values: [][]driver.Value{{"greptime_value", "double", ""}, {"greptime_timestamp", "timestamp", "TIME INDEX"}}},
		"FROM information_schema.partitions":     {columns: []string{"partition_name", "partition_ordinal_position", "partition_expression", "partition_description", "greptime_partition_id"}, values: [][]driver.Value{{"p0", int64(1), "region", "region = 'east'", int64(42)}}},
		"is_leader = 'Yes'":                      {columns: []string{"region_id", "peer_id"}, values: [][]driver.Value{{int64(42), "datanode-1"}}},
	})
	result, err := Discover(context.Background(), db, []string{"metrics"})
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(result.Tables) != 1 || len(result.Skipped) != 0 {
		t.Fatalf("tables=%d skipped=%#v", len(result.Tables), result.Skipped)
	}
	table := result.Tables[0]
	if table.PhysicalTable != "physical" || table.TimeIndex != "greptime_timestamp" || table.Partitions[0].RegionID != 42 || table.Partitions[0].LeaderDatanode != "datanode-1" {
		t.Fatalf("unexpected discovery: %#v", table)
	}
	if table.Partitions[0].Predicate.SQL != "`region` = ?" {
		t.Fatalf("unexpected predicate: %#v", table.Partitions[0].Predicate)
	}
}

func TestDiscoverUsesPartitionMetadataAsPhysicalPartitionAuthority(t *testing.T) {
	db := mockDiscoveryDB(t, map[string]mockRows{
		"FROM information_schema.tables":         {columns: []string{"table_name"}, values: [][]driver.Value{{"cpu"}}},
		"SHOW CREATE TABLE `metrics`.`cpu`":      {columns: []string{"Table", "Create Table"}, values: [][]driver.Value{{"cpu", "CREATE TABLE cpu WITH (on_physical_table = 'physical')"}}},
		"SHOW CREATE TABLE `metrics`.`physical`": {columns: []string{"Table", "Create Table"}, values: [][]driver.Value{{"physical", "CREATE TABLE physical PARTITION ON COLUMNS (region) (region = 'east')"}}},
		"FROM information_schema.columns":        {columns: []string{"column_name", "data_type", "column_key"}, values: [][]driver.Value{{"value", "double", ""}, {"ts", "timestamp", "TIME INDEX"}}},
		"FROM information_schema.partitions":     {columns: []string{"partition_name", "partition_ordinal_position", "partition_expression", "partition_description", "greptime_partition_id"}, values: [][]driver.Value{{"p-west", int64(3), "region", "region = 'west'", int64(99)}}},
		"is_leader = 'Yes'":                      {columns: []string{"region_id", "peer_id"}, values: [][]driver.Value{{int64(99), "datanode-2"}}},
	})
	result, err := Discover(context.Background(), db, []string{"metrics"})
	if err != nil || len(result.Tables) != 1 || result.Tables[0].Partitions[0].Name != "p-west" || result.Tables[0].Partitions[0].RegionID != 99 {
		t.Fatalf("result=%#v err=%v", result, err)
	}
}

func TestDiscoverDoesNotQueryAmbiguousPartitionRows(t *testing.T) {
	db := mockDiscoveryDB(t, map[string]mockRows{
		"FROM information_schema.tables":         {columns: []string{"table_name"}, values: [][]driver.Value{{"cpu"}}},
		"SHOW CREATE TABLE `metrics`.`cpu`":      {columns: []string{"Table", "Create Table"}, values: [][]driver.Value{{"cpu", "CREATE TABLE cpu WITH (on_physical_table = 'physical')"}}},
		"SHOW CREATE TABLE `metrics`.`physical`": {columns: []string{"Table", "Create Table"}, values: [][]driver.Value{{"physical", "CREATE TABLE physical PARTITION ON COLUMNS (region) (region = 'east')"}}},
		"FROM information_schema.columns":        {columns: []string{"column_name", "data_type", "column_key"}, values: [][]driver.Value{{"value", "double", ""}, {"ts", "timestamp", "TIME INDEX"}}},
		"FROM information_schema.partitions":     {columns: []string{"partition_name", "partition_ordinal_position", "partition_expression", "partition_description", "greptime_partition_id"}, values: [][]driver.Value{{"duplicated", int64(1), "region", "region = 'east'", int64(1)}, {"duplicated", int64(2), "region", "region = 'west'", int64(2)}, {"p2", int64(3), "region", "region = 'central'", int64(3)}}},
		"is_leader = 'Yes'":                      {columns: []string{"region_id", "peer_id"}, values: [][]driver.Value{{int64(1), "datanode-3"}, {int64(2), "datanode-3"}, {int64(3), "datanode-3"}}},
	})
	result, err := Discover(context.Background(), db, []string{"metrics"})
	if err != nil || len(result.Tables) != 1 || len(result.Tables[0].Partitions) != 1 || result.Tables[0].Partitions[0].Name != "p2" {
		t.Fatalf("result=%#v err=%v", result, err)
	}
}

func TestRegionLeaderAcceptsNullPeerID(t *testing.T) {
	db := mockDiscoveryDB(t, map[string]mockRows{
		"is_leader = 'Yes'": {columns: []string{"peer_id"}, values: [][]driver.Value{{nil}}},
	})
	leader, err := regionLeader(context.Background(), db, 42)
	if err != nil || leader != "" {
		t.Fatalf("leader=%q err=%v", leader, err)
	}
}

func TestDiscoverSkipsMissingMetadataAndUnsafeTables(t *testing.T) {
	db := mockDiscoveryDB(t, map[string]mockRows{
		"FROM information_schema.tables":         {columns: []string{"table_name"}, values: [][]driver.Value{{"missing"}, {"unsafe"}}},
		"SHOW CREATE TABLE `metrics`.`missing`":  {columns: []string{"Table", "Create Table"}, values: [][]driver.Value{{"missing", "CREATE TABLE missing WITH (on_physical_table = 'physical') PARTITION ON COLUMNS (region) (region = 'east')"}}},
		"SHOW CREATE TABLE `metrics`.`physical`": {columns: []string{"Table", "Create Table"}, values: [][]driver.Value{{"physical", "CREATE TABLE physical PARTITION ON COLUMNS (region) (region = 'east')"}}},
		"SHOW CREATE TABLE `metrics`.`unsafe`":   {columns: []string{"Table", "Create Table"}, values: [][]driver.Value{{"unsafe", "CREATE TABLE unsafe"}}},
		"FROM information_schema.columns":        {columns: []string{"column_name", "data_type", "column_key"}, values: [][]driver.Value{{"value", "double", ""}, {"ts", "timestamp", "TIME INDEX"}}},
		"FROM information_schema.partitions":     {columns: []string{"partition_name", "partition_ordinal_position", "partition_expression", "partition_description", "greptime_partition_id"}},
	})
	result, err := Discover(context.Background(), db, []string{"metrics"})
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(result.Tables) != 0 || len(result.Skipped) != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if !strings.Contains(result.Skipped[0].Reason, "partition metadata") || !strings.Contains(result.Skipped[1].Reason, "on_physical_table") {
		t.Fatalf("unexpected skip reasons: %#v", result.Skipped)
	}
}

type mockRows struct {
	columns []string
	values  [][]driver.Value
}
type discoveryDriver struct{ responses map[string]mockRows }
type discoveryConn struct{ responses map[string]mockRows }
type discoveryRows struct {
	mockRows
	index int
}

func (d discoveryDriver) Open(string) (driver.Conn, error)  { return discoveryConn(d), nil }
func (c discoveryConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (c discoveryConn) Close() error                        { return nil }
func (c discoveryConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }
func (c discoveryConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	for key, rows := range c.responses {
		if strings.Contains(query, key) {
			return &discoveryRows{mockRows: rows}, nil
		}
	}
	return nil, fmt.Errorf("unexpected query: %s", query)
}
func (r *discoveryRows) Columns() []string { return r.columns }
func (r *discoveryRows) Close() error      { return nil }
func (r *discoveryRows) Next(dest []driver.Value) error {
	if r.index == len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.index])
	r.index++
	return nil
}
func mockDiscoveryDB(t *testing.T, responses map[string]mockRows) *sql.DB {
	t.Helper()
	name := "discovery_mock_" + strings.ReplaceAll(t.Name(), "/", "_")
	sql.Register(name, discoveryDriver{responses})
	db, err := sql.Open(name, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
