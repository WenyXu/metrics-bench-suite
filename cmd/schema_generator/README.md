# Schema Generator

SchemaGenerator is a tool to generate SQL schema for physical and metrics tables, along with sample loader configurations. It can also optionally execute the generated SQL against a MySQL-compatible database.

## Usage

```bash
./bin/schema_generator [flags]
```

**Example:**

Generate SQL schemas and sample loader YAML files, then execute the SQLs against a local GreptimeDB instance:

```bash
./bin/schema_generator \
  --physical-table-count 5 \
  --column-list-count 50 \
  --selected-column-count-range "3,7" \
  --sample-count 3 \
  --region-count 10 \
  --target-path ./generated_schema \
  --exec-sql \
  --mysql-host 127.0.0.1 \
  --mysql-port 4002 \
  --do-exec-sql-job-count 2
```

## Flags

*   `-H, --mysql-host string`: The MySQL host (default: "127.0.0.1"), GreptimeDB instance address.
*   `-P, --mysql-port string`: The MySQL port (default: "4002"), GreptimeDB instance mysql port.
*   `-r, --region-count int`: The number of regions to generate for every physical table (default: 50)
*   `-p, --physical-table-count int`: The number of physical tables to generate (default: 100)
*   `-c, --column-list-count int`: The number of columns in the physical table (default: 100)
*   `-s, --selected-column-count-range string`: The range of selected column count for metrics tables, e.g., "5,10" (default: "5,10"), which means the number of columns in each metrics table will be randomly selected from this range.
*   `-n, --sample-count int`: The number of sample metric tables to generate per physical table (default: 10)
*   `-t, --target-path string`: The target path to save the generated files (default: "yaml_config")
*   `-d, --exec-sql`: Execute SQL files after generation (default: false)
*   `-j, --do-exec-sql-job-count uint`: The number of jobs to execute SQL concurrently (default: 4)

## Output

The tool generates the following files in the specified `--target-path`:

1.  `physical_table.sql`: Contains `CREATE TABLE` statements for the physical tables.
2.  `metrics_table.sql`: Contains `CREATE TABLE` statements for the metrics tables.
3.  `metrics_table_info.json`: A JSON file containing information about the generated metrics tables, including their names, associated columns, and the physical table they are based on.
4.  YAML files (e.g., `metrics_table_0.yaml`, `metrics_table_1.yaml`, ...): Sample configuration files for the `sample_loader` tool, one for each generated metrics table. These files define how to generate sample data for the corresponding table.

If `--exec-sql` is true, the SQL statements in `physical_table.sql` and `metrics_table.sql` will be executed against the specified MySQL-compatible database.