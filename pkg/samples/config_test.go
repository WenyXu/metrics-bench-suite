package samples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	yamlConfig := `
start: "2023-05-09T17:49:12+08:00"
end: "2023-05-09T18:49:12+08:00"
interval: 10 # 10s
precision: 1000 # ms
num-field: 10
tags:
  - name: hostname
    type: STRING
    nullability: 0
    dist:
      type: constant_string
      value: host-0
  - name: region
    type: STRING
    nullability: 0
    dist:
      type: weighted_preset
      preset:
        - value: sg-central
          weight: 1
        - value: us-east-1
          weight: 1
        - value: us-west-2
          weight: 1
        - value: us-central-1
          weight: 1
  - name: datacenter
    type: STRING
    nullability: 0
    dist:
      type: weighted_preset
      preset:
        - value: dc-0
          weight: 1
        - value: dc-1
          weight: 1
  - name: rack
    type: STRING
    nullability: 0
    dist:
      type: weighted_preset
      preset:
        - value: "1"
          weight: 1
        - value: "2"
          weight: 1
        - value: "3"
          weight: 1
        - value: "4"
          weight: 1
  - name: os
    type: STRING
    nullability: 0
    dist:
      type: weighted_preset
      preset:
        - value: ubuntu-22.04
          weight: 1
        - value: ubuntu-20.04
          weight: 0.5
  - name: arch
    type: STRING
    nullability: 0
    dist:
      type: weighted_preset
      preset:
        - value: x86_64
          weight: 1
        - value: aarch64
          weight: 1
  - name: team
    type: STRING
    nullability: 0
    dist:
      type: weighted_preset
      preset:
        - value: HZ
          weight: 15
        - value: BJ
          weight: 10
        - value: SH
          weight: 0
        - value: GZ
          weight: 2
  - name: services
    type: STRING
    nullability: 0
    dist:
      type: weighted_preset
      preset:
        - value: etcd
          weight: 1
        - value: kafka
          weight: 1
        - value: frontend
          weight: 10
        - value: datanode
          weight: 5
        - value: metasrv
          weight: 3
  - name: services_version
    type: STRING
    nullability: 0
    dist:
      type: weighted_preset
      preset:
        - value: "1.0"
          weight: 1
        - value: "2.0"
          weight: 1
  - name: services_env
    type: STRING
    nullability: 0
    dist:
      type: weighted_preset
      preset:
        - value: offline
          weight: 1
        - value: test
          weight: 2
        - value: prod
          weight: 1

fields:
  - name: cpu
    type: FLOAT
    nullability: 0
    dist:
      type: random
      lower_bound: 0
      upper_bound: 100
  - name: mem
    type: FLOAT
    nullability: 0
    dist:
      type: normal
      mean: 0.5
      stddev: 0.1
  - name: disk
    type: FLOAT
    nullability: 0
    dist:
      type: mono_inc
      step: 11
  - name: network
    type: FLOAT
    nullability: 0
    dist:
      type: normal
      mean: 100
      stddev: 2000
  - name: io
    type: FLOAT
    nullability: 0
    dist:
      type: normal
      mean: 0
      stddev: 200
  - name: cpu_temp
    type: FLOAT
    nullability: 0
    dist:
      type: uniform
      lower_bound: 0
      upper_bound: 100
  - name: mem_temp
    type: FLOAT
    nullability: 0
    dist:
      type: noise
      max_fluctuation: 30
  - name: disk_temp
    type: FLOAT
    nullability: 0
    dist:
      type: normal
      mean: 0.5
      stddev: 0.1
  - name: network_temp
    type: FLOAT
    nullability: 0
    dist:
      type: periodic
      period: 100
      amplitude: 204800
      bias: 1024
  - name: io_temp
    type: FLOAT
    nullability: 0
    dist:
      type: constant_float
      value: 10.1
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlConfig), &config)
	if err != nil {
		t.Fatalf("failed to unmarshal yaml config: %v", err)
	}

	assert.Equal(t, config.Start, "2023-05-09T17:49:12+08:00")
	assert.Equal(t, config.End, "2023-05-09T18:49:12+08:00")
	assert.Equal(t, config.Interval, 10)
	assert.Equal(t, config.Tags[0].Name, "hostname")
	assert.Equal(t, config.Tags[0].Type, "STRING")
	assert.Equal(t, config.Tags[0].Dist.Type, "constant_string")
	assert.Equal(t, config.Tags[0].Dist.Value.(string), "host-0")
	assert.Equal(t, config.Tags[1].Name, "region")
	assert.Equal(t, config.Tags[1].Type, "STRING")
	assert.Equal(t, config.Tags[1].Dist.Type, "weighted_preset")
	assert.Equal(t, config.Tags[1].Dist.Preset[0].Value, "sg-central")
	assert.Equal(t, config.Tags[1].Dist.Preset[0].Weight, 1)
	assert.Equal(t, config.Tags[1].Dist.Preset[1].Value, "us-east-1")
	assert.Equal(t, config.Tags[1].Dist.Preset[1].Weight, 1)
	assert.Equal(t, config.Tags[1].Dist.Preset[2].Value, "us-west-2")
	assert.Equal(t, config.Tags[1].Dist.Preset[2].Weight, 1)
	assert.Equal(t, config.Tags[1].Dist.Preset[3].Value, "us-central-1")
}
