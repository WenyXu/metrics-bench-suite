### Sample Loader

The sample loader is a tool genrate samples and send them to the remote write endpoint.

#### Usage

Generate and load the sample data from the config file.

```bash
./bin/sample_loader -c ./configs/debug_samples_400 -u  http://localhost:4000/v1/prometheus/write\?db\=public --start-date 2025-03-09T18:00:00+08:00 --end-date 2025-03-09T19:00:00+08:00 --interval 30s  --tick-interval 1s
```

#### Configs
`./configs/debug_samples` total time series 4k.
`./configs/debug_samples_400` total time series 400k.
`./configs/debug_samples_800` total time series 800k.
