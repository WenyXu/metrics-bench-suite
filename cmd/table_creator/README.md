## Table Creator

The table creator is a tool that creates a table sql for the physical table and logical tables.

### Usage

It will write the sql to the file `metrics1.greptime_physical_table-create-tables.sql`.
```bash
./bin/table_creator -c ./configs/debug_samples/  -d metrics1
```

### Flags

- `-c`: The path to the config files.
- `-d`: The name of the database.
- `-t`: The name of the physical table.
- `-s`: The skipping index granularity.



