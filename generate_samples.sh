#!/bin/bash

START_TIME=2025-03-09T00:00:00Z
END_TIME=2025-03-09T00:20:00Z
FILE_DURATION=10m

# Convert FILE_DURATION to seconds
DURATION_SECS=$(echo $FILE_DURATION | sed 's/m/*60/' | bc)

CURRENT_TIME=$START_TIME
COUNTER=1

while [[ "$CURRENT_TIME" < "$END_TIME" ]]; do
    # Calculate next end time
    NEXT_TIME=$(date -u -d "$CURRENT_TIME + $DURATION_SECS seconds" --iso-8601=seconds | sed 's/+00:00/Z/')
    
    if [[ "$NEXT_TIME" > "$END_TIME" ]]; then
        NEXT_TIME=$END_TIME
    fi
    
    echo "Generating data from $CURRENT_TIME to $NEXT_TIME"
    ./bin/sample_generator -c ./configs/debug_samples_1500 --start-date "$CURRENT_TIME" --end-date "$NEXT_TIME" --interval 30s 
    
    CURRENT_TIME=$NEXT_TIME
    ((COUNTER++))
done
