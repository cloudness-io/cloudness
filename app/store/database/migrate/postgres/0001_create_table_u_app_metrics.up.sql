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