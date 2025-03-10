#!/bin/bash

# Usage: ./import_sql.sh <num_dbs>
# Example: ./import_sql.sh 5

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <num_dbs>"
    exit 1
fi

NUM_DBS=$1

for i in $(seq 1 $NUM_DBS); do
    DB_NAME="metrics$i"
    echo "Importing SQL files for database: $DB_NAME"
    mysql -h 127.0.0.1 -P 4002 < "$DB_NAME.greptime_physical_table-create-tables.sql"
done