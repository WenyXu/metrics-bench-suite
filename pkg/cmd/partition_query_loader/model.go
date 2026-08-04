package partition_query_loader

// DiscoveredTable is a logical metric table that can safely receive bounded reads.
type DiscoveredTable struct {
	Database      string
	LogicalTable  string
	PhysicalTable string
	ValueColumn   string
	TimeIndex     string
	Partitions    []DiscoveredPartition
}

// DiscoveredPartition joins a parsed partition predicate with its region placement.
type DiscoveredPartition struct {
	Name           string
	Description    string
	RegionID       uint64
	LeaderDatanode string
	DatanodeID     string
	Predicate      Predicate
}

// DiscoveryResult contains usable tables and the reasons other tables were skipped.
type DiscoveryResult struct {
	Tables  []DiscoveredTable
	Skipped []SkipReason
}

// SkipReason records why discovery declined to create a potentially unsafe query.
type SkipReason struct {
	Database string
	Table    string
	Reason   string
}
