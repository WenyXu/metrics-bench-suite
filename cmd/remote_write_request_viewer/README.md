# Remote Write Request Viewer

Remote Write Request Viewer is a command-line tool designed to decode Prometheus remote write body data (snappy compressed).

## Usage

The tool can be used from the command line. Below are some examples of how to use it:

### Decode from Input
```bash
remote_write_request_viewer <base64_encoded_data>
```

### Decode from File (Base64 Encoded)

```bash
remote_write_request_viewer -t file <path/to/file>
```

## Example
```bash
remote_write_request_viewer U/BSClEKHwoIX19uYW1lX18SE2h0dHBfcmVxdWVzdHNfdG90YWwKDQoGc3RhdHVzEgMyMDAKDQoGbWV0aG9kEgNHRVQSEAkAAAAAAADwPxCMuO6ZyDI=
```
output:
```bash
Decoded WriteRequest:

Time Series 1:
Labels:
  __name__ = http_requests_total
  status = 200
  method = GET
Samples:
  Value: 1.000000, Timestamp: 1737368509452
```