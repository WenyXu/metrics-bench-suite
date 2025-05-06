CREATE DATABASE IF NOT EXISTS `label_inverted_200`;

CREATE TABLE `label_inverted_200`.`greptime_physical_table` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`id` STRING NULL INVERTED INDEX,
`kubernetes_namespace` STRING NULL INVERTED INDEX,
`priority_class` STRING NULL INVERTED INDEX,
`method` STRING NULL INVERTED INDEX,
`plugin_name` STRING NULL INVERTED INDEX,
`node` STRING NULL INVERTED INDEX,
`code` STRING NULL INVERTED INDEX,
`host_ip` STRING NULL INVERTED INDEX,
`state` STRING NULL INVERTED INDEX,
`env` STRING NULL INVERTED INDEX,
`resource` STRING NULL INVERTED INDEX,
`volumemode` STRING NULL INVERTED INDEX,
`owner_is_controller` STRING NULL INVERTED INDEX,
`host` STRING NULL INVERTED INDEX,
`prometheus` STRING NULL INVERTED INDEX,
`storageclass` STRING NULL INVERTED INDEX,
`pod_ip` STRING NULL INVERTED INDEX,
`owner_kind` STRING NULL INVERTED INDEX,
`unit` STRING NULL INVERTED INDEX,
`interface` STRING NULL INVERTED INDEX,
`reason` STRING NULL INVERTED INDEX,
`volumename` STRING NULL INVERTED INDEX,
`owner_name` STRING NULL INVERTED INDEX,
`workload_type` STRING NULL INVERTED INDEX,
`volume_plugin` STRING NULL INVERTED INDEX,
`plugin` STRING NULL INVERTED INDEX,
`app_kubernetes_io_instance` STRING NULL INVERTED INDEX,
`app_kubernetes_io_managed_by` STRING NULL INVERTED INDEX,
`migrated` STRING NULL INVERTED INDEX,
`uid` STRING NULL INVERTED INDEX,
`name` STRING NULL INVERTED INDEX,
`app_kubernetes_io_component` STRING NULL INVERTED INDEX,
`persistentvolumeclaim` STRING NULL INVERTED INDEX,
`metrics_path` STRING NULL INVERTED INDEX,
`rcode` STRING NULL INVERTED INDEX,
`instance_type` STRING NULL INVERTED INDEX,
`kubernetes_pod_name` STRING NULL INVERTED INDEX,
`holiday` STRING NULL INVERTED INDEX,
`namespace` STRING NULL INVERTED INDEX,
`image` STRING NULL INVERTED INDEX,
`zone` STRING NULL INVERTED INDEX,
`app_kubernetes_io_part_of` STRING NULL INVERTED INDEX,
`phase` STRING NULL INVERTED INDEX,
`kubernetes_name` STRING NULL INVERTED INDEX,
`instance` STRING NULL INVERTED INDEX,
`service` STRING NULL INVERTED INDEX,
`cpu` STRING NULL INVERTED INDEX,
`app` STRING NULL INVERTED INDEX,
`proto` STRING NULL INVERTED INDEX,
`created_by_kind` STRING NULL INVERTED INDEX,
`cluster` STRING NULL INVERTED INDEX,
`verb` STRING NULL INVERTED INDEX,
`type` STRING NULL INVERTED INDEX,
`condition` STRING NULL INVERTED INDEX,
`operation_type` STRING NULL INVERTED INDEX,
`argocd_argoproj_io_instance` STRING NULL INVERTED INDEX,
`device` STRING NULL INVERTED INDEX,
`app_kubernetes_io_version` STRING NULL INVERTED INDEX,
`helm_sh_chart` STRING NULL INVERTED INDEX,
`hostname` STRING NULL INVERTED INDEX,
`region` STRING NULL INVERTED INDEX,
`container` STRING NULL INVERTED INDEX,
`endpoint` STRING NULL INVERTED INDEX,
`job` STRING NULL INVERTED INDEX,
`container_state` STRING NULL INVERTED INDEX,
`workload` STRING NULL INVERTED INDEX,
`status` STRING NULL INVERTED INDEX,
`created_by_name` STRING NULL INVERTED INDEX,
`host_network` STRING NULL INVERTED INDEX,
`operation_name` STRING NULL INVERTED INDEX,
`pod` STRING NULL INVERTED INDEX,
`le` STRING NULL INVERTED INDEX,
`family` STRING NULL INVERTED INDEX,
`app_kubernetes_io_name` STRING NULL INVERTED INDEX,
`prometheus_replica` STRING NULL INVERTED INDEX,
`server` STRING NULL INVERTED INDEX,
`zones` STRING NULL INVERTED INDEX,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`id`,`kubernetes_namespace`,`priority_class`,`method`,`plugin_name`,`node`,`code`,`host_ip`,`state`,`env`,`resource`,`volumemode`,`owner_is_controller`,`host`,`prometheus`,`storageclass`,`pod_ip`,`owner_kind`,`unit`,`interface`,`reason`,`volumename`,`owner_name`,`workload_type`,`volume_plugin`,`plugin`,`app_kubernetes_io_instance`,`app_kubernetes_io_managed_by`,`migrated`,`uid`,`name`,`app_kubernetes_io_component`,`persistentvolumeclaim`,`metrics_path`,`rcode`,`instance_type`,`kubernetes_pod_name`,`holiday`,`namespace`,`image`,`zone`,`app_kubernetes_io_part_of`,`phase`,`kubernetes_name`,`instance`,`service`,`cpu`,`app`,`proto`,`created_by_kind`,`cluster`,`verb`,`type`,`condition`,`operation_type`,`argocd_argoproj_io_instance`,`device`,`app_kubernetes_io_version`,`helm_sh_chart`,`hostname`,`region`,`container`,`endpoint`,`job`,`container_state`,`workload`,`status`,`created_by_name`,`host_network`,`operation_name`,`pod`,`le`,`family`,`app_kubernetes_io_name`,`prometheus_replica`,`server`,`zones`,)
) ENGINE = metric WITH (
    "physical_metric_table" = "",   
    "memtable.type" = "partition_tree",
    "memtable.partition_tree.primary_key_encoding" = "sparse",
	"index.type" = "inverted", 
	"index.granularity" = "102400",
	"compaction.type" = "twcs",
	"compaction.twcs.time_window" = "2h",
);

CREATE TABLE `label_inverted_200`.`apiserver_request:availability30d` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`env` STRING,
`holiday` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`verb` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`env`,`holiday`,`prometheus`,`prometheus_replica`,`verb`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`cluster:namespace:pod_cpu:active:kube_pod_container_resource_limits` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`resource` STRING,
`service` STRING,
`uid` STRING,
`unit` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`resource`,`service`,`uid`,`unit`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`cluster:namespace:pod_cpu:active:kube_pod_container_resource_requests` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`resource` STRING,
`service` STRING,
`uid` STRING,
`unit` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`resource`,`service`,`uid`,`unit`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`cluster:namespace:pod_memory:active:kube_pod_container_resource_limits` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`resource` STRING,
`service` STRING,
`uid` STRING,
`unit` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`resource`,`service`,`uid`,`unit`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`cluster:namespace:pod_memory:active:kube_pod_container_resource_requests` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`resource` STRING,
`service` STRING,
`uid` STRING,
`unit` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`resource`,`service`,`uid`,`unit`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`cluster:node_cpu:ratio_rate5m` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`env` STRING,
`holiday` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`env`,`holiday`,`prometheus`,`prometheus_replica`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`cluster_quantile:apiserver_request_slo_duration_seconds:histogram_quantile` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`code_resource:apiserver_request_total:rate5m` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`code` STRING,
`env` STRING,
`holiday` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`resource` STRING,
`verb` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`code`,`env`,`holiday`,`prometheus`,`prometheus_replica`,`resource`,`verb`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_cpu_cfs_periods_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_cpu_cfs_throttled_periods_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_cpu_usage_seconds_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`cpu` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`cpu`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_fs_limit_bytes` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_fs_reads_bytes_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`device` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`device`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_fs_reads_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`device` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`device`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_fs_usage_bytes` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_fs_writes_bytes_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`device` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`device`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_fs_writes_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`device` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`device`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_memory_cache` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_memory_rss` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_memory_swap` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_memory_working_set_bytes` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_network_receive_bytes_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`interface` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`interface`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_network_receive_packets_dropped_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`interface` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`interface`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_network_receive_packets_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`interface` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`interface`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_network_transmit_bytes_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`interface` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`interface`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_network_transmit_packets_dropped_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`interface` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`interface`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`container_network_transmit_packets_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`interface` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`interface`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_cache_entries` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`server` STRING,
`service` STRING,
`type` STRING,
`zones` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`server`,`service`,`type`,`zones`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_cache_hits_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`server` STRING,
`service` STRING,
`type` STRING,
`zones` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`server`,`service`,`type`,`zones`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_cache_misses_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`server` STRING,
`service` STRING,
`zones` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`server`,`service`,`zones`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_cache_size` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_dns_do_requests_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_dns_request_count_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_dns_request_do_count_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_dns_request_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`server` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`server`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_dns_request_size_bytes_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`proto` STRING,
`server` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`proto`,`server`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_dns_request_type_count_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_dns_requests_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`family` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`proto` STRING,
`server` STRING,
`service` STRING,
`type` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`family`,`holiday`,`instance`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`proto`,`server`,`service`,`type`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_dns_response_rcode_count_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_dns_response_size_bytes_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`proto` STRING,
`server` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`proto`,`server`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`coredns_dns_responses_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`plugin` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`rcode` STRING,
`server` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`plugin`,`pod`,`prometheus`,`prometheus_replica`,`rcode`,`server`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`go_goroutines` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`app` STRING,
`app_kubernetes_io_component` STRING,
`app_kubernetes_io_instance` STRING,
`app_kubernetes_io_managed_by` STRING,
`app_kubernetes_io_name` STRING,
`app_kubernetes_io_part_of` STRING,
`app_kubernetes_io_version` STRING,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`helm_sh_chart` STRING,
`holiday` STRING,
`hostname` STRING,
`instance` STRING,
`instance_type` STRING,
`job` STRING,
`kubernetes_namespace` STRING,
`kubernetes_pod_name` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`region` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`app`,`app_kubernetes_io_component`,`app_kubernetes_io_instance`,`app_kubernetes_io_managed_by`,`app_kubernetes_io_name`,`app_kubernetes_io_part_of`,`app_kubernetes_io_version`,`cluster`,`container`,`endpoint`,`env`,`helm_sh_chart`,`holiday`,`hostname`,`instance`,`instance_type`,`job`,`kubernetes_namespace`,`kubernetes_pod_name`,`metrics_path`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`region`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_namespace_status_phase` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`namespace` STRING,
`job` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`namespace`,`job`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_node_info` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`node` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`node`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_node_status_allocatable` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`resource` STRING,
`service` STRING,
`unit` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`resource`,`service`,`unit`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_node_status_capacity` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`resource` STRING,
`service` STRING,
`unit` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`resource`,`service`,`unit`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_node_status_condition` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`condition` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
`status` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`condition`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,`status`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_persistentvolumeclaim_info` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`persistentvolumeclaim` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
`storageclass` STRING,
`volumemode` STRING,
`volumename` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`persistentvolumeclaim`,`pod`,`prometheus`,`prometheus_replica`,`service`,`storageclass`,`volumemode`,`volumename`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_persistentvolumeclaim_status_phase` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`persistentvolumeclaim` STRING,
`phase` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`persistentvolumeclaim`,`phase`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_pod_container_resource_limits` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`resource` STRING,
`service` STRING,
`uid` STRING,
`unit` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`resource`,`service`,`uid`,`unit`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_pod_container_resource_requests` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`resource` STRING,
`service` STRING,
`uid` STRING,
`unit` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`resource`,`service`,`uid`,`unit`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_pod_container_status_last_terminated_reason` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`reason` STRING,
`service` STRING,
`uid` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`reason`,`service`,`uid`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_pod_info` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`created_by_kind` STRING,
`created_by_name` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`host_ip` STRING,
`host_network` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`pod_ip` STRING,
`priority_class` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
`uid` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`created_by_kind`,`created_by_name`,`endpoint`,`env`,`holiday`,`host_ip`,`host_network`,`instance`,`job`,`namespace`,`node`,`pod`,`pod_ip`,`priority_class`,`prometheus`,`prometheus_replica`,`service`,`uid`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_pod_owner` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`owner_is_controller` STRING,
`owner_kind` STRING,
`owner_name` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
`uid` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`owner_is_controller`,`owner_kind`,`owner_name`,`pod`,`prometheus`,`prometheus_replica`,`service`,`uid`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kube_resourcequota` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_cgroup_manager_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`operation_type` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`metrics_path`,`namespace`,`node`,`operation_type`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_cgroup_manager_duration_seconds_count` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`operation_type` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`operation_type`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_node_config_error` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_node_name` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_pleg_relist_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`metrics_path`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_pleg_relist_duration_seconds_count` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_pleg_relist_interval_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`metrics_path`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_pod_start_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`metrics_path`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_pod_start_duration_seconds_count` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_pod_worker_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`operation_type` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`metrics_path`,`namespace`,`node`,`operation_type`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_pod_worker_duration_seconds_count` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`operation_type` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`operation_type`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_running_container_count` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_running_containers` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container_state` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container_state`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_running_pod_count` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_running_pods` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_runtime_operations_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`operation_type` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`metrics_path`,`namespace`,`node`,`operation_type`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_runtime_operations_errors_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`operation_type` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`operation_type`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_runtime_operations_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`operation_type` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`operation_type`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_volume_stats_available_bytes` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`persistentvolumeclaim` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`persistentvolumeclaim`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_volume_stats_capacity_bytes` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`persistentvolumeclaim` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`persistentvolumeclaim`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubelet_volume_stats_used_bytes` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`persistentvolumeclaim` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`persistentvolumeclaim`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubeproxy_network_programming_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubeproxy_network_programming_duration_seconds_count` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubeproxy_sync_proxy_rules_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`kubeproxy_sync_proxy_rules_duration_seconds_count` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`namespace_cpu:kube_pod_container_resource_limits:sum` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`env` STRING,
`holiday` STRING,
`namespace` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`env`,`holiday`,`namespace`,`prometheus`,`prometheus_replica`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`namespace_cpu:kube_pod_container_resource_requests:sum` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`env` STRING,
`holiday` STRING,
`namespace` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`env`,`holiday`,`namespace`,`prometheus`,`prometheus_replica`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`namespace_memory:kube_pod_container_resource_limits:sum` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`env` STRING,
`holiday` STRING,
`namespace` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`env`,`holiday`,`namespace`,`prometheus`,`prometheus_replica`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`namespace_memory:kube_pod_container_resource_requests:sum` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`env` STRING,
`holiday` STRING,
`namespace` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`env`,`holiday`,`namespace`,`prometheus`,`prometheus_replica`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`namespace_workload_pod:kube_pod_owner:relabel` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`env` STRING,
`holiday` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`workload` STRING,
`workload_type` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`env`,`holiday`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`workload`,`workload_type`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`node_memory_MemAvailable_bytes:sum` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`node_memory_MemTotal_bytes` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`hostname` STRING,
`instance` STRING,
`instance_type` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`region` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`hostname`,`instance`,`instance_type`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`region`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`env` STRING,
`holiday` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`env`,`holiday`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`node_namespace_pod_container:container_memory_cache` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`node_namespace_pod_container:container_memory_rss` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`node_namespace_pod_container:container_memory_swap` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`node_namespace_pod_container:container_memory_working_set_bytes` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`id` STRING,
`image` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`id`,`image`,`instance`,`job`,`metrics_path`,`name`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`node_netstat_TcpExt_TCPSynRetrans` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`hostname` STRING,
`instance` STRING,
`instance_type` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`region` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`hostname`,`instance`,`instance_type`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`region`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`node_netstat_Tcp_OutSegs` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`hostname` STRING,
`instance` STRING,
`instance_type` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`region` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`hostname`,`instance`,`instance_type`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`region`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`node_netstat_Tcp_RetransSegs` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`hostname` STRING,
`instance` STRING,
`instance_type` STRING,
`job` STRING,
`namespace` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`region` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`container`,`endpoint`,`env`,`holiday`,`hostname`,`instance`,`instance_type`,`job`,`namespace`,`pod`,`prometheus`,`prometheus_replica`,`region`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`process_cpu_seconds_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`app` STRING,
`app_kubernetes_io_component` STRING,
`app_kubernetes_io_managed_by` STRING,
`app_kubernetes_io_name` STRING,
`app_kubernetes_io_part_of` STRING,
`app_kubernetes_io_version` STRING,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`hostname` STRING,
`instance` STRING,
`instance_type` STRING,
`job` STRING,
`kubernetes_name` STRING,
`kubernetes_namespace` STRING,
`kubernetes_pod_name` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`region` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`app`,`app_kubernetes_io_component`,`app_kubernetes_io_managed_by`,`app_kubernetes_io_name`,`app_kubernetes_io_part_of`,`app_kubernetes_io_version`,`cluster`,`container`,`endpoint`,`env`,`holiday`,`hostname`,`instance`,`instance_type`,`job`,`kubernetes_name`,`kubernetes_namespace`,`kubernetes_pod_name`,`metrics_path`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`region`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`process_resident_memory_bytes` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`app` STRING,
`app_kubernetes_io_component` STRING,
`app_kubernetes_io_instance` STRING,
`app_kubernetes_io_managed_by` STRING,
`app_kubernetes_io_name` STRING,
`app_kubernetes_io_part_of` STRING,
`app_kubernetes_io_version` STRING,
`argocd_argoproj_io_instance` STRING,
`cluster` STRING,
`container` STRING,
`endpoint` STRING,
`env` STRING,
`helm_sh_chart` STRING,
`holiday` STRING,
`hostname` STRING,
`instance` STRING,
`instance_type` STRING,
`job` STRING,
`kubernetes_name` STRING,
`kubernetes_namespace` STRING,
`kubernetes_pod_name` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`region` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`app`,`app_kubernetes_io_component`,`app_kubernetes_io_instance`,`app_kubernetes_io_managed_by`,`app_kubernetes_io_name`,`app_kubernetes_io_part_of`,`app_kubernetes_io_version`,`argocd_argoproj_io_instance`,`cluster`,`container`,`endpoint`,`env`,`helm_sh_chart`,`holiday`,`hostname`,`instance`,`instance_type`,`job`,`kubernetes_name`,`kubernetes_namespace`,`kubernetes_pod_name`,`metrics_path`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`region`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`rest_client_request_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`host` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
`verb` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`host`,`instance`,`job`,`le`,`metrics_path`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,`verb`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`rest_client_requests_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`code` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`host` STRING,
`instance` STRING,
`job` STRING,
`kubernetes_namespace` STRING,
`kubernetes_pod_name` STRING,
`method` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`code`,`endpoint`,`env`,`holiday`,`host`,`instance`,`job`,`kubernetes_namespace`,`kubernetes_pod_name`,`method`,`metrics_path`,`namespace`,`node`,`pod`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`storage_operation_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`le` STRING,
`metrics_path` STRING,
`migrated` STRING,
`namespace` STRING,
`node` STRING,
`operation_name` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
`status` STRING,
`volume_plugin` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`le`,`metrics_path`,`migrated`,`namespace`,`node`,`operation_name`,`prometheus`,`prometheus_replica`,`service`,`status`,`volume_plugin`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`storage_operation_duration_seconds_count` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`migrated` STRING,
`namespace` STRING,
`node` STRING,
`operation_name` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
`status` STRING,
`volume_plugin` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`migrated`,`namespace`,`node`,`operation_name`,`prometheus`,`prometheus_replica`,`service`,`status`,`volume_plugin`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`storage_operation_errors_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
TIME INDEX (greptime_timestamp),
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`up` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`app` STRING,
`app_kubernetes_io_component` STRING,
`app_kubernetes_io_managed_by` STRING,
`app_kubernetes_io_name` STRING,
`app_kubernetes_io_part_of` STRING,
`app_kubernetes_io_version` STRING,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`hostname` STRING,
`instance` STRING,
`instance_type` STRING,
`job` STRING,
`kubernetes_name` STRING,
`kubernetes_namespace` STRING,
`kubernetes_pod_name` STRING,
`metrics_path` STRING,
`pod` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`region` STRING,
`service` STRING,
`zone` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`app`,`app_kubernetes_io_component`,`app_kubernetes_io_managed_by`,`app_kubernetes_io_name`,`app_kubernetes_io_part_of`,`app_kubernetes_io_version`,`cluster`,`endpoint`,`env`,`holiday`,`hostname`,`instance`,`instance_type`,`job`,`kubernetes_name`,`kubernetes_namespace`,`kubernetes_pod_name`,`metrics_path`,`pod`,`prometheus`,`prometheus_replica`,`region`,`service`,`zone`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`volume_manager_total_volumes` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`metrics_path` STRING,
`namespace` STRING,
`node` STRING,
`plugin_name` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
`state` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`metrics_path`,`namespace`,`node`,`plugin_name`,`prometheus`,`prometheus_replica`,`service`,`state`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`workqueue_adds_total` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`kubernetes_namespace` STRING,
`kubernetes_pod_name` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`kubernetes_namespace`,`kubernetes_pod_name`,`metrics_path`,`name`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`workqueue_depth` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`kubernetes_namespace` STRING,
`kubernetes_pod_name` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`kubernetes_namespace`,`kubernetes_pod_name`,`metrics_path`,`name`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);

CREATE TABLE `label_inverted_200`.`workqueue_queue_duration_seconds_bucket` (
greptime_timestamp TIMESTAMP(3) NOT NULL,
greptime_value DOUBLE NULL,
`cluster` STRING,
`endpoint` STRING,
`env` STRING,
`holiday` STRING,
`instance` STRING,
`job` STRING,
`kubernetes_namespace` STRING,
`kubernetes_pod_name` STRING,
`le` STRING,
`metrics_path` STRING,
`name` STRING,
`namespace` STRING,
`node` STRING,
`prometheus` STRING,
`prometheus_replica` STRING,
`service` STRING,
TIME INDEX (greptime_timestamp),
PRIMARY KEY (`cluster`,`endpoint`,`env`,`holiday`,`instance`,`job`,`kubernetes_namespace`,`kubernetes_pod_name`,`le`,`metrics_path`,`name`,`namespace`,`node`,`prometheus`,`prometheus_replica`,`service`,)
) ENGINE = metric WITH (
on_physical_table = 'greptime_physical_table'
);
