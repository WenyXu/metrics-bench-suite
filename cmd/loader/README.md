# Metrics Bench Suite Loader

The Metrics Bench Suite Loader is designed to load time series data, where the time series labels are derived from traffic captured in a real-world environment (via tcpdump), along with randomly generated values and timestamps, into the target database.

## Features

- Load data into a database using Prometheus write requests.
- Supports dry-run mode for testing without actual data processing.
- Configurable time series per request and interval.
- Uses a random seed for generating data.


### Flags

- `-u, --url`: The URL of the database (default: `http://localhost:4000/v1/prometheus/write?db=public`).
- `-t, --tcpflow-output`: The path to the tcpflow output.
- `-r, --timeseries-per-request`: The number of timeseries per request(remote write request) (default: `2000`).
- `-i, --interval`: The interval of the loading data (default: `10s`).
- `-s, --seed`: The seed for the random number generator (default: `123456`).
- `--start-date`: The start date of the data (default: `2025-01-01T00:00:00Z`).
- `-d, --dry-run`: Dry run the loader without processing data.


## Usage

### Dump tcp packets 
To capture remote write requests from a real-world environment, such as VictoriaMetrics or any other backend that supports Prometheus remote write requests.

```bash
tcpdump -i eth0 tcp port <port> -w prometheus_requests.pcap
```

### Reassemble the packets
```bash
tcpflow -r prometheus_requests.pcap -o ./tcpflow_output
```

### Load data into the GrpeitmeDB
```bash
loader -t <tcpflow_output_file> -url http://localhost:4000/v1/prometheus/write?db=public
```

## Example

To run the loader with a specific configuration:

```bash
./bin/loader --url "http://your-database-url" -t "path/to/tcpflow-output" -r 1000 -i "5s" -s 654321 -d "2023-01-01T00:00:00Z" 
```
