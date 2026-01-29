package types

import "time"

type AppMetrics struct {
	Timestamp      time.Time `db:"matric_timestamp"              json:"timestamp"`
	ApplicationUID int64     `db:"matric_application_uid"        json:"application_uid"`
	InstanceName   string    `db:"matric_instance_name"          json:"instance_name"`
	CPU            int64     `db:"matric_cpu_mcores"             json:"cpu"`
	Memory         int64     `db:"matric_memory_bytes"           json:"memory"`
}

type AppMetricsAggregate struct {
	BucketTimestamp int64   `db:"bucket_ts"             json:"bucket_ts"`
	InstanceName    string  `db:"metric_instance_name"  json:"instance_name"`
	CPU             float64 `db:"cpu_vcores"            json:"cpu"`
	Memory          float64 `db:"memory_bytes"          json:"memory"`
}

type AppMetricsViewModel struct {
	CPU    []*MetricsSeriesViewModel
	Memory []*MetricsSeriesViewModel
}

type MetricsSeriesViewModel struct {
	Label      string    `json:"label"`
	Timestamps []int64   `json:"timestamps"`
	Values     []float64 `json:"values"`
}
