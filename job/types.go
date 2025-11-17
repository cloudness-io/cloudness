package job

type Job struct {
	UID                 string   `db:"job_uid"`
	Created             int64    `db:"job_created"`
	Updated             int64    `db:"job_updated"`
	Type                string   `db:"job_type"`
	Priority            Priority `db:"job_priority"`
	Data                string   `db:"job_data"`
	Result              string   `db:"job_result"`
	MaxDurationSeconds  int      `db:"job_max_duration_seconds"`
	MaxRetries          int      `db:"job_max_retries"`
	State               State    `db:"job_state"`
	Scheduled           int64    `db:"job_scheduled"`
	TotalExecutions     int      `db:"job_total_executions"`
	RunBy               string   `db:"job_run_by"`
	RunDeadline         int64    `db:"job_run_deadline"`
	RunProgress         int      `db:"job_run_progress"`
	LastExecuted        int64    `db:"job_last_executed"`
	IsRecurring         bool     `db:"job_is_recurring"`
	RecurringCron       string   `db:"job_recurring_cron"`
	ConsecutiveFailures int      `db:"job_consecutive_failures"`
	LastFailureError    string   `db:"job_last_failure_error"`
	GroupID             string   `db:"job_group_id"`
}

type StateChange struct {
	UID      string `json:"uid"`
	Type     string `json:"type"`
	State    State  `json:"state"`
	Progress int    `json:"progress"`
	Result   string `json:"result"`
	Failure  string `json:"failure"`
}

type Progress struct {
	State    State  `json:"state"`
	Progress int    `json:"progress"`
	Result   string `json:"result,omitempty"`
	Failure  string `json:"failure,omitempty"`
}
