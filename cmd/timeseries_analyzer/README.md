# Timeseries Analyzer

The Timeseries Analyzer analyzes tcpdump data and outputs the number of timeseries in remote write requests.

## Usage

### Dump tcp packets
```bash
tcpdump -i eth0 tcp port 4000 -w prometheus_requests.pcap
```

### Reassemble the packets
```bash
tcpflow -r prometheus_requests.pcap -o ./tcpflow_output
```

### Analyze the tcpflow output
```bash
timeseries_analyzer <tcpflow_output_file>
```

## Example
```bash
timeseries_analyzer ./tcpflow_output/010.000.010.148.52136-010.000.059.168.04000
```
output:
```bash
2025/01/20 16:44:46 Processing file: ./tcpflow_output/010.000.010.148.52136-010.000.059.168.04000
2025/01/20 16:44:46 Found 2087 table names
2025/01/20 16:44:46 Found 74998 time series
```