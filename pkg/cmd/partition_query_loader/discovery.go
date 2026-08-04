package partition_query_loader

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var physicalTablePattern = regexp.MustCompile(`(?is)\bon_physical_table\s*=\s*'((?:''|[^'])+)'`)
var tableIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_:]*$`)

// DiscoverConfigured opens the configured MySQL endpoint and discovers its eligible tables.
func DiscoverConfigured(ctx context.Context, config Config) (DiscoveryResult, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, config.DiscoveryTimeout)
	defer cancel()
	db, err := config.OpenDB()
	if err != nil {
		return DiscoveryResult{}, err
	}
	defer db.Close()
	return DiscoverWithOptions(ctx, db, config.Databases, DiscoveryOptions{TablesPerDatabase: config.TablesPerDatabase, RandomSeed: config.RandomSeed})
}

// Discover queries metadata for logical metric tables and returns only tables with safe bounded predicates.
func Discover(ctx context.Context, db *sql.DB, databases []string) (DiscoveryResult, error) {
	return discover(ctx, db, databases, DiscoveryOptions{})
}

// DiscoveryOptions bounds table selection and reports progress without affecting discovery.
type DiscoveryOptions struct {
	TablesPerDatabase int
	RandomSeed        *int64
	Progress          func(database string, selected, total, completed int)
}

// DiscoverWithOptions discovers tables using the supplied bounded-selection options.
func DiscoverWithOptions(ctx context.Context, db *sql.DB, databases []string, options DiscoveryOptions) (DiscoveryResult, error) {
	return discover(ctx, db, databases, options)
}

func discover(ctx context.Context, db *sql.DB, databases []string, options DiscoveryOptions) (DiscoveryResult, error) {
	result := DiscoveryResult{}
	cache := discoveryCache{physical: make(map[string]*physicalDiscovery)}
	var random *rand.Rand
	if options.RandomSeed != nil {
		random = rand.New(rand.NewSource(*options.RandomSeed))
	} else {
		random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	for _, database := range databases {
		database = strings.TrimSpace(database)
		if !identifierPattern.MatchString(database) {
			return DiscoveryResult{}, fmt.Errorf("invalid database identifier %q", database)
		}
		tables, err := logicalTables(ctx, db, database)
		if err != nil {
			return DiscoveryResult{}, err
		}
		totalTables := len(tables)
		tables = selectTables(tables, options.TablesPerDatabase, random)
		if options.Progress != nil {
			options.Progress(database, len(tables), totalTables, 0)
		}
		for completed, table := range tables {
			discovered, reason := discoverTable(ctx, db, database, table, &cache)
			if options.Progress != nil {
				options.Progress(database, len(tables), totalTables, completed+1)
			}
			if reason != "" {
				result.Skipped = append(result.Skipped, SkipReason{Database: database, Table: table, Reason: reason})
				continue
			}
			result.Tables = append(result.Tables, discovered)
		}
	}
	return result, nil
}

func selectTables(tables []string, limit int, random *rand.Rand) []string {
	selected := append([]string(nil), tables...)
	if limit <= 0 || len(selected) <= limit {
		return selected
	}
	if random == nil {
		random = rand.New(rand.NewSource(1))
	}
	for i := len(selected) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		selected[i], selected[j] = selected[j], selected[i]
	}
	return selected[:limit]
}

func logicalTables(ctx context.Context, db *sql.DB, database string) ([]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_type = 'BASE TABLE'`, database)
	if err != nil {
		return nil, fmt.Errorf("list tables in %s: %w", database, err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		if tableIdentifierPattern.MatchString(table) {
			tables = append(tables, table)
		}
	}
	return tables, rows.Err()
}

type physicalDiscovery struct {
	definition PartitionDefinition
	metadata   []partitionMetadataRow
	leaders    map[uint64]string
	err        string
}

type discoveryCache struct{ physical map[string]*physicalDiscovery }

func discoverTable(ctx context.Context, db *sql.DB, database, table string, cache *discoveryCache) (DiscoveredTable, string) {
	showCreate, err := showCreateTable(ctx, db, database, table)
	if err != nil {
		return DiscoveredTable{}, fmt.Sprintf("read table definition: %v", err)
	}
	physicalTable, ok := physicalTable(showCreate)
	if !ok {
		return DiscoveredTable{}, "not a logical metric table with on_physical_table"
	}
	valueColumn, timeIndex, err := queryColumns(ctx, db, database, table)
	if err != nil {
		return DiscoveredTable{}, err.Error()
	}
	key := database + "\x00" + physicalTable
	physical := cache.physical[key]
	if physical == nil {
		physical = &physicalDiscovery{}
		physicalCreate, readErr := showCreateTable(ctx, db, database, physicalTable)
		if readErr != nil {
			physical.err = fmt.Sprintf("read physical table definition: %v", readErr)
		} else if physical.definition, readErr = ParsePartitionDefinition(physicalCreate); readErr != nil {
			physical.err = readErr.Error()
		}
		if physical.err == "" {
			physical.metadata, readErr = partitionMetadata(ctx, db, database, physicalTable)
			if readErr != nil {
				physical.err = readErr.Error()
			}
		}
		if physical.err == "" {
			physical.leaders, readErr = regionLeaders(ctx, db, physical.metadata)
			if readErr != nil {
				physical.err = readErr.Error()
			}
		}
		cache.physical[key] = physical
	}
	if physical.err != "" {
		return DiscoveredTable{}, physical.err
	}
	metadata := physical.metadata
	if len(metadata) == 0 {
		return DiscoveredTable{}, "physical table has no partition metadata"
	}

	predicates := metadataPredicates(physical.definition, metadata)
	partitions := make([]DiscoveredPartition, 0, len(predicates))
	for index, predicate := range predicates {
		if predicate == nil {
			continue
		}
		item := metadata[index]
		leader := physical.leaders[item.regionID]
		if leader == "" {
			return DiscoveredTable{}, fmt.Sprintf("region %d has no leader", item.regionID)
		}
		partitions = append(partitions, DiscoveredPartition{Name: item.name, Description: item.description, RegionID: item.regionID, LeaderDatanode: leader, Predicate: *predicate})
	}
	if len(partitions) == 0 {
		return DiscoveredTable{}, "no partition metadata row has a safe predicate"
	}
	return DiscoveredTable{Database: database, LogicalTable: table, PhysicalTable: physicalTable, ValueColumn: valueColumn, TimeIndex: timeIndex, Partitions: partitions}, ""
}

// metadataPredicates derives predicates from authoritative PARTITIONS rows. It only
// consults the physical DDL for rows whose description cannot be parsed directly.
func metadataPredicates(definition PartitionDefinition, metadata []partitionMetadataRow) []*Predicate {
	predicates := make([]*Predicate, len(metadata))
	uniqueNames, uniqueOrdinals := metadataIdentities(metadata)
	for index, item := range metadata {
		if !uniqueNames[item.name] || !uniqueOrdinals[item.ordinal] {
			continue
		}
		if predicate, ok := metadataPredicate(item); ok {
			predicates[index] = &predicate
		}
	}
	for index, item := range metadata {
		if predicates[index] != nil || !uniqueNames[item.name] || !uniqueOrdinals[item.ordinal] || item.ordinal < 1 || item.ordinal > uint64(len(definition.Partitions)) {
			continue
		}
		columns, err := parseColumns(item.expression)
		if err != nil || !sameColumns(columns, definition.Columns) {
			continue
		}
		predicate, err := BuildPartitionPredicate(definition.Partitions[item.ordinal-1])
		if err == nil {
			predicates[index] = &predicate
		}
	}
	return predicates
}

func metadataIdentities(metadata []partitionMetadataRow) (map[string]bool, map[uint64]bool) {
	nameCounts := make(map[string]int, len(metadata))
	ordinalCounts := make(map[uint64]int, len(metadata))
	for _, item := range metadata {
		nameCounts[item.name]++
		ordinalCounts[item.ordinal]++
	}
	names := make(map[string]bool, len(nameCounts))
	for name, count := range nameCounts {
		names[name] = name != "" && count == 1
	}
	ordinals := make(map[uint64]bool, len(ordinalCounts))
	for ordinal, count := range ordinalCounts {
		ordinals[ordinal] = ordinal != 0 && count == 1
	}
	return names, ordinals
}

func metadataPredicate(item partitionMetadataRow) (Predicate, bool) {
	columns, err := parseColumns(item.expression)
	if err != nil {
		return Predicate{}, false
	}
	partition, err := parsePartition(item.description, columns)
	if err != nil {
		return Predicate{}, false
	}
	predicate, err := BuildPartitionPredicate(partition)
	return predicate, err == nil
}

func sameColumns(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func showCreateTable(ctx context.Context, db *sql.DB, database, table string) (string, error) {
	rows, err := db.QueryContext(ctx, "SHOW CREATE TABLE "+quoteIdentifier(database)+"."+quoteIdentifier(table))
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if !rows.Next() {
		return "", fmt.Errorf("no CREATE TABLE result")
	}
	var name, create string
	if err := rows.Scan(&name, &create); err != nil {
		return "", err
	}
	return create, rows.Err()
}

func physicalTable(showCreate string) (string, bool) {
	matches := physicalTablePattern.FindStringSubmatch(showCreate)
	if len(matches) != 2 {
		return "", false
	}
	name := strings.ReplaceAll(matches[1], "''", "'")
	if !tableIdentifierPattern.MatchString(name) {
		return "", false
	}
	return name, true
}

func queryColumns(ctx context.Context, db *sql.DB, database, table string) (string, string, error) {
	rows, err := db.QueryContext(ctx, "SELECT column_name, data_type, column_key FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ordinal_position", database, table)
	if err != nil {
		return "", "", fmt.Errorf("read table columns: %w", err)
	}
	defer rows.Close()
	value, timeIndex := "", ""
	for rows.Next() {
		var name, dataType, columnKey string
		if err := rows.Scan(&name, &dataType, &columnKey); err != nil {
			return "", "", err
		}
		if strings.EqualFold(columnKey, "TIME INDEX") && identifierPattern.MatchString(name) {
			timeIndex = name
		}
		if value == "" && isNumericType(dataType) && identifierPattern.MatchString(name) {
			value = name
		}
	}
	if err := rows.Err(); err != nil {
		return "", "", err
	}
	if value == "" || timeIndex == "" {
		return "", "", fmt.Errorf("missing numeric aggregate column or time index")
	}
	return value, timeIndex, nil
}

func isNumericType(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "tinyint", "smallint", "int", "integer", "bigint", "float", "double", "real", "decimal":
		return true
	default:
		return false
	}
}

type partitionMetadataRow struct {
	name, description string
	regionID          uint64
	expression        string
	ordinal           uint64
}

func partitionMetadata(ctx context.Context, db *sql.DB, database, physicalTable string) ([]partitionMetadataRow, error) {
	rows, err := db.QueryContext(ctx, "SELECT partition_name, partition_ordinal_position, partition_expression, partition_description, greptime_partition_id FROM information_schema.partitions WHERE table_schema = ? AND table_name = ? ORDER BY partition_ordinal_position", database, physicalTable)
	if err != nil {
		return nil, fmt.Errorf("read partition metadata: %w", err)
	}
	defer rows.Close()
	var metadata []partitionMetadataRow
	for rows.Next() {
		var item partitionMetadataRow
		if err := rows.Scan(&item.name, &item.ordinal, &item.expression, &item.description, &item.regionID); err != nil {
			return nil, err
		}
		metadata = append(metadata, item)
	}
	return metadata, rows.Err()
}

func regionLeader(ctx context.Context, db *sql.DB, regionID uint64) (string, error) {
	rows, err := db.QueryContext(ctx, "SELECT peer_id FROM information_schema.region_peers WHERE region_id = ? AND is_leader = 'Yes'", regionID)
	if err != nil {
		return "", fmt.Errorf("read leader for region %d: %w", regionID, err)
	}
	defer rows.Close()
	if !rows.Next() {
		return "", rows.Err()
	}
	var leader sql.NullString
	if err := rows.Scan(&leader); err != nil {
		return "", err
	}
	if !leader.Valid {
		return "", rows.Err()
	}
	return leader.String, rows.Err()
}

func regionLeaders(ctx context.Context, db *sql.DB, metadata []partitionMetadataRow) (map[uint64]string, error) {
	ids := make([]uint64, 0, len(metadata))
	seen := make(map[uint64]bool)
	for _, item := range metadata {
		if !seen[item.regionID] {
			seen[item.regionID] = true
			ids = append(ids, item.regionID)
		}
	}
	leaders := make(map[uint64]string, len(ids))
	if len(ids) == 0 {
		return leaders, nil
	}
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	rows, err := db.QueryContext(ctx, "SELECT region_id, peer_id FROM information_schema.region_peers WHERE is_leader = 'Yes' AND region_id IN ("+strings.Join(placeholders, ",")+")", args...)
	if err != nil {
		return nil, fmt.Errorf("read leaders for physical table: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id uint64
		var leader sql.NullString
		if err := rows.Scan(&id, &leader); err != nil {
			return nil, err
		}
		if leader.Valid {
			leaders[id] = leader.String
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return leaders, nil
}
