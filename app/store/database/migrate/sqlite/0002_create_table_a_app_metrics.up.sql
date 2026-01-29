CREATE TABLE app_metrics (
 metric_timestamp               INTEGER NOT NULL
,metric_application_uid 		INTEGER NOT NULL
,metric_instance_name           TEXT NOT NULL
,metric_cpu_mcores              INTEGER NOT NULL
,metric_memory_bytes            INTEGER NOT NULL

,UNIQUE(metric_timestamp, metric_application_uid, metric_instance_name)
);

CREATE INDEX app_metrics_timestamp
ON app_metrics(metric_timestamp, metric_application_uid, metric_instance_name);