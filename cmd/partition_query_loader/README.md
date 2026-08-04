# Partition Query Loader / 分区查询负载工具

`partition_query_loader` discovers GreptimeDB logical metric tables, maps their physical-table partitions to regions and leaders, then issues bounded aggregate reads against the logical tables. It connects through the MySQL protocol on port `4002`.

## Install / 安装

```bash
go run ./cmd/partition_query_loader --help
```

The command uses the repository's MySQL driver and needs no additional service-side installation.

## Safety / 安全限制

The loader only uses `SELECT`, `SHOW CREATE TABLE`, and `information_schema` metadata queries. It never sends DDL or DML. Every workload query requires a parsed partition predicate and a bounded time range; unsupported partition syntax and tables without safe query columns are skipped. Workload execution is enabled by default; use `--dry-run` to inspect plans without issuing reads. Initialization and one-time scheduled-plan logs, plus per-sample CPU-rate logs, go to stderr; NDJSON/CSV output is unchanged and requests are not logged individually.

## Flags / 参数

- `--mysql-host`, `--mysql-port` (default `4002`), `--mysql-user`, `--mysql-password`: MySQL connection settings.
- `--databases metrics1,metrics2`: databases to discover.
- `--tables-per-database N`: uniformly select at most `N` eligible logical tables in each database before per-table metadata discovery; `0` (default) selects all. `--random-seed` makes selection reproducible.
- `--dry-run` (default `false`): execute reads; otherwise only print plans and skipped reasons.
- `--profile sustained|periodic`, `--duration` (default `0`, unlimited until cancelled): request schedule. Every plan has an independent fixed cadence and shares only the global concurrency limit; hot/cold settings deterministically skew those cadences.
- `--period-min`, `--period-max`: inclusive per-plan execution interval range, sampled once at startup (defaults `15s`–`45s`). Deprecated `--period` sets both bounds.
- `--concurrency`, `--hot-partitions`, `--hot-share`: concurrency and hotspot skew.
- `--time-range-min`, `--time-range-max`: inclusive per-plan fixed window length, sampled once at startup (defaults `5m`–`15m`; hot plans use the upper half and cold plans the lower half). Deprecated `--time-range` sets both bounds. Each execution refreshes the latest timestamp under its partition predicate; empty partitions are skipped.
- `--hot-partitions`, `--hot-share`: number of hot plans and target request-rate share. Hot plans are selected preferentially from the largest initially co-located datanode group; the configured seed makes ties and selection reproducible.
- `--stats-interval`, `--output-format ndjson|csv`: region-statistics sampling and final output encoding. Samples refresh `datanode_id` and log per-datanode region counts and valid total CPU cores.
- `--discovery-timeout`, `--mysql-dial-timeout`, `--mysql-read-timeout`, `--mysql-write-timeout`: bounded discovery and MySQL operations.

## Examples / 示例

Inspect safe plans first:

```bash
go run ./cmd/partition_query_loader --mysql-host 127.0.0.1 --mysql-user root --databases metrics1 --dry-run
```

Create sustained skewed reads and save NDJSON observations:

```bash
go run ./cmd/partition_query_loader --mysql-host 127.0.0.1 --mysql-user root --databases metrics1 --dry-run=false --profile sustained --duration 5m --period 30s --hot-partitions 3 --hot-share 0.85 > observations.ndjson
```

Use sparse periodic reads and CSV output:

```bash
go run ./cmd/partition_query_loader --mysql-host 127.0.0.1 --mysql-user root --databases metrics1 --dry-run=false --profile periodic --duration 20m --period 5m --stats-interval 15s --output-format csv > observations.csv
```

## Output / 输出

NDJSON and CSV observations include table and partition identity, `region_id`, `leader_datanode`, cumulative CPU and scanned-byte counters, CPU-core and scan-byte rates, request/error counts, p95 latency, plus `moved` and `reset` flags. A leader change or a disappeared region is marked `moved`; decreasing counters are marked `reset` and never produce negative rates.

## Manual MySQL-4002 verification / 手工验证

1. Start a GreptimeDB cluster with the MySQL endpoint reachable on `4002`.
2. Run the dry-run example and confirm that every printed query has both a partition predicate and timestamp bounds.
3. Run the sustained example and compare `query_cpu_time_rate_cores` with the region-balancer or Grafana read-load trend.
4. Track the same `region_id`; a changed `leader_datanode` confirms a migration. Continue sampling after migration to assess whether load became more balanced.
5. Run the periodic example to evaluate estimator behavior for sparse reads.
