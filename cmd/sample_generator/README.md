## Sample Generator

The sample generator is a tool that generates sample data for a given config file.

### Usage

#### Generate samples and save to file
```bash
./bin/sample_generator -c ./assets -o output.json
```

#### Generate samples and send to remote write

**No output file will be generated.**
```bash
./bin/sample_generator -c ./assets -u http://localhost:4000/v1/prometheus/write?db=public --start-date 2024-03-07T00:00:00+08:00 --end-date 2024-03-08T00:00:00+08:00 --interval 30s
```


### Config File

The config file is a yaml file that describes the data to be generated. **The filename will be used as the metric name**

The config file is in the `./assets` directory.

For example, `./assets/kubelet_node_name.yaml`
```yaml
tags:
  - name: env
    type: STRING
    dist:
      type: constant_string
      value: 'prod'
  - name: holiday
    type: String
    dist:
      type: constant_string
      value: 'false'
  - name: instance
    type: String
    dist:
      type: replica_string
      replica_prefix: '10.0.162.'
      replica: 30
  - name: job_name
    type: String
    dist:
      type: constant_string
      value: 'kubelet'
  - name: metrics_path
    type: String
    dist:
      type: constant_string
      value: '/metrics'
  - name: namespace
    type: String
    dist:
      type: constant_string
      value: 'kube-system'
fields:
  - name: greptime_value
    type: Float
    dist:
      type: periodic
      period: 1000
      amplitude: 100
      bias: 100
```

### Flags
- `-c`: The path to the config files.
- `-o`: The output file.
- `-u`: The remote write url.
- `-i`: The interval of the data to be generated.
- `-s`: The seed of the random number generator.
- `-start-date`: The start date of the data to be generated.
- `-end-date`: The end date of the data to be generated.

