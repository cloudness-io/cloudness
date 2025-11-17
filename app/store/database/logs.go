package database

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"

	"github.com/jmoiron/sqlx"
)

var _ store.LogStore = (*LogStore)(nil)

type LogStore struct {
	db *sqlx.DB
}

func NewLogStore(db *sqlx.DB) *LogStore {
	return &LogStore{
		db: db,
	}
}

type logs struct {
	DeploymentID int64  `db:"log_deployment_id"`
	Data         []byte `db:"log_data"`
}

const logsColumns = `
   log_deployment_id
   ,log_data
`

const logsInsert = `
INSERT INTO logs (
	log_deployment_id
	,log_data
) VALUES (
	:log_deployment_id
	,:log_data
)
`

const logsBaseQuery = `
SELECT` + logsColumns + `
FROM logs`

func (s *LogStore) Find(ctx context.Context, deploymentID int64) (io.ReadCloser, error) {
	const sqlQuery = logsBaseQuery + ` WHERE log_deployment_id = $1`
	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(logs)
	if err := db.GetContext(ctx, dst, sqlQuery, deploymentID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to find log")
	}

	return io.NopCloser(bytes.NewReader(dst.Data)), nil
}

func (s *LogStore) Create(ctx context.Context, deploymentID int64, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not read log: %w", err)
	}

	params := &logs{
		DeploymentID: deploymentID,
		Data:         data,
	}

	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(logsInsert, params)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to bind log")
	}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "log query insert failed")
	}

	return nil
}
