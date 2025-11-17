package types

// LogLine represents a line in the logs.
type LogLine struct {
	Number    int    `json:"pos"`
	Message   string `json:"out"`
	Timestamp int64  `json:"time"`
}

// Logs represets the full logs of a deployment.
type Logs struct {
	ID           int64  `db:"log_id"            json:"id"`
	DeploymentID int64  `db:"log_deployment_id" json:"deployment_id"`
	Data         []byte `db:"log_data"          json:"data"`
}
