 CREATE EXTENSION IF NOT EXISTS timescaledb;

 CREATE TABLE app_metrics (
  metric_timestamp        TIMESTAMPTZ       NOT NULL,
  metric_application_uid  BIGINT            NOT NULL,
  metric_instance_name    TEXT              NOT NULL,
  metric_cpu_mcores       INTEGER           NOT NULL,
  metric_memory_bytes     BIGINT            NOT NULL,

  PRIMARY KEY (metric_timestamp, metric_application_uid, metric_instance_name)
);

SELECT create_hypertable(
  'app_metrics',
  'metric_timestamp',
  chunk_time_interval => INTERVAL '1 day',
  if_not_exists => true
);

CREATE INDEX app_metrics_app_instance_time
ON app_metrics (metric_application_uid, metric_instance_name, metric_timestamp DESC);

ALTER TABLE app_metrics
SET (
  timescaledb.compress,
  timescaledb.compress_segmentby = 'metric_application_uid, metric_instance_name'
);

SELECT add_compression_policy(
  'app_metrics',
  INTERVAL '7 days'
);

SELECT add_retention_policy(
  'app_metrics',
  INTERVAL '30 days'
);