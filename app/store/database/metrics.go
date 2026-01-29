package database

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	"github.com/jmoiron/sqlx"
)

var _ store.MetricsStore = (*MetricsStore)(nil)

// NewMetricsStore returns a new metrics store
func NewMetricsStore(db *sqlx.DB) *MetricsStore {
	return &MetricsStore{db: db}
}

// MetricsStore defines the metrics data storage
type MetricsStore struct {
	db *sqlx.DB
}

// metricsColumns defines the columns that are same across all metrics.
const metricsColumns = `
	metric_timestamp
	,metric_application_uid
	,metric_instance_name
	,metric_cpu_mcores
	,metric_memory_bytes`

// UpsertMany upserts the metrics
func (s *MetricsStore) UpsertMany(ctx context.Context, metrics []*types.AppMetrics) error {
	driver := s.db.DriverName()

	stmt := database.Builder.Insert("app_metrics").
		Columns(`metric_timestamp
					,metric_application_uid
					,metric_instance_name
					,metric_cpu_mcores
					,metric_memory_bytes`)

	for _, m := range metrics {
		var timestamp interface{}

		switch driver {
		case "sqlite3":
			// SQLite: Convert time.Time to ISO8601 string
			timestamp = m.Timestamp.UTC().Format("2006-01-02 15:04:05")
		default:
			// PostgreSQL/TimescaleDB: Use time.Time directly
			timestamp = m.Timestamp
		}

		stmt = stmt.Values(
			timestamp,
			m.ApplicationUID,
			m.InstanceName,
			m.CPU,
			m.Memory,
		)
	}

	stmt = stmt.Suffix(`ON CONFLICT (metric_timestamp, metric_application_uid, metric_instance_name) 
	DO UPDATE SET 
		metric_cpu_mcores = EXCLUDED.metric_cpu_mcores
		,metric_memory_bytes = EXCLUDED.metric_memory_bytes`)

	query, args, err := stmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to convert upsert template query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "metrics upsert query failed")
	}

	return nil
}

// ListByApplicationUID retrieves aggregated metrics for an application within a time range
func (s *MetricsStore) ListByApplicationUID(
	ctx context.Context,
	applicationUID int64,
	fromTime time.Time,
	toTime time.Time,
	bucketSeconds int64,
) ([]*types.AppMetricsAggregate, error) {

	driver := s.db.DriverName()

	var query string
	var args []any

	switch driver {
	case "sqlite3":
		// SQLite: Use strftime for time bucketing
		query = `
		SELECT
			strftime('%s', metric_timestamp) / ? * ? AS bucket_ts,
			metric_instance_name,
			AVG(metric_cpu_mcores) / 1000.0 AS cpu_vcores,
			AVG(metric_memory_bytes) AS memory_bytes
		FROM app_metrics
		WHERE metric_application_uid = ?
		  AND metric_timestamp BETWEEN ? AND ?
		GROUP BY bucket_ts, metric_instance_name
		ORDER BY bucket_ts ASC
		`
		args = []any{
			bucketSeconds,
			bucketSeconds,
			applicationUID,
			fromTime.UTC().Format("2006-01-02 15:04:05"),
			toTime.UTC().Format("2006-01-02 15:04:05"),
		}

	default: // postgres / timescale
		query = `
		SELECT
			(EXTRACT(EPOCH FROM metric_timestamp)::BIGINT / $1) * $1 AS bucket_ts,
			metric_instance_name,
			AVG(metric_cpu_mcores) / 1000.0 AS cpu_vcores,
			AVG(metric_memory_bytes) AS memory_bytes
		FROM app_metrics
		WHERE metric_application_uid = $2
		  AND metric_timestamp BETWEEN $3 AND $4
		GROUP BY bucket_ts, metric_instance_name
		ORDER BY bucket_ts ASC
		`
		args = []any{
			bucketSeconds,
			applicationUID,
			fromTime,
			toTime,
		}
	}

	db := dbtx.GetAccessor(ctx, s.db)
	rows := []*types.AppMetricsAggregate{}
	if err := db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by application uid query failed")
	}

	return rows, nil
}

// DeleteOldMetrics removes metrics older than the specified duration
func (s *MetricsStore) DeleteOldMetrics(ctx context.Context, olderThan time.Duration) (int64, error) {
	driver := s.db.DriverName()
	cutoffTime := time.Now().Add(-olderThan)

	var query string
	var args []any

	switch driver {
	case "sqlite3":
		query = `
		DELETE FROM app_metrics
		WHERE metric_timestamp < ?
		`
		args = []any{cutoffTime.UTC().Format("2006-01-02 15:04:05")}

	default: // postgres / timescale
		query = `
		DELETE FROM app_metrics
		WHERE metric_timestamp < $1
		`
		args = []any{cutoffTime}
	}

	db := dbtx.GetAccessor(ctx, s.db)
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Delete old metrics query failed")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Failed to get rows affected")
	}

	// Run VACUUM on SQLite if we deleted a lot of rows
	if driver == "sqlite3" && rowsAffected > 1000 {
		if _, err := db.ExecContext(ctx, "VACUUM"); err != nil {
			// Log but don't fail - VACUUM is optimization
		}
	}

	return rowsAffected, nil
}
