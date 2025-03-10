#!/bin/bash

# Usage: ./generate_databases.sh <num_dbs>
# Example: ./generate_databases.sh 5

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <num_dbs>"
    exit 1
fi

NUM_DBS=$1

for i in $(seq 1 $NUM_DBS); do
    DB_NAME="metrics$i"
    echo "Creating database: $DB_NAME"
    ./bin/table_creator -c ./configs/debug_samples/ -d "$DB_NAME"
done