package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	baseStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.DeploymentStore = (*DeploymentStore)(nil)

// DeploymentStore implements a DeploymentStore backed by a relational database.
type DeploymentStore struct {
	db *sqlx.DB
}

// NewDeploymentStore returns new DeploymentStore
func NewDeploymentStore(db *sqlx.DB) *DeploymentStore {
	return &DeploymentStore{
		db: db,
	}
}

// deployment is a DB representation of a deployment.
type deployment struct {
	types.Deployment
}

// deploymentColumns defines the columns that are the same across all deployments.
const deploymentColumns = `
   deployment_id
	,deployment_uid
   ,deployment_application_id
	,deployment_spec
	,deployment_needs_build
   ,deployment_triggerer
   ,deployment_title
	,deployment_action
	,deployment_status
   ,deployment_error
	,deployment_version
	,deployment_machine
   ,deployment_started
	,deployment_stopped
   ,deployment_created
   ,deployment_updated`

const deploymentInsert = `
INSERT INTO deployments  (
   deployment_uid
   ,deployment_application_id
   ,deployment_spec
   ,deployment_needs_build
   ,deployment_triggerer
   ,deployment_title
	,deployment_action
   ,deployment_status
   ,deployment_error
   ,deployment_version
	,deployment_machine
   ,deployment_started
   ,deployment_stopped
   ,deployment_created
   ,deployment_updated
) values (
   :deployment_uid
   ,:deployment_application_id
   ,:deployment_spec
   ,:deployment_needs_build
	,:deployment_triggerer
	,:deployment_title
	,:deployment_action
   ,:deployment_status
   ,:deployment_error
   ,:deployment_version
	,:deployment_machine
   ,:deployment_started
   ,:deployment_stopped
   ,:deployment_created
   ,:deployment_updated
	) RETURNING deployment_id
	`

const deploymentBaseQuery = `
SELECT` + deploymentColumns + `
FROM deployments`

func (s *DeploymentStore) Create(ctx context.Context, deployment *types.Deployment) (*types.Deployment, error) {
	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(deploymentInsert, deployment)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind deployment object")
	}

	if err := db.QueryRowContext(ctx, query, arg...).Scan(&deployment.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert deployment query failed")
	}
	return deployment, nil
}

func (s *DeploymentStore) Find(ctx context.Context, id int64) (*types.Deployment, error) {
	const sqlQuery = deploymentBaseQuery + ` WHERE deployment_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(deployment)

	if err := db.GetContext(ctx, dst, sqlQuery, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}

	return s.mapDBDeployment(dst)
}

func (s *DeploymentStore) FindByUID(ctx context.Context, applicationID int64, uid int64) (*types.Deployment, error) {
	const sqlQuery = deploymentBaseQuery + ` WHERE deployment_application_id = $1 AND deployment_uid = $2`

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(deployment)

	if err := db.GetContext(ctx, dst, sqlQuery, applicationID, uid); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by uid query failed")
	}

	return s.mapDBDeployment(dst)
}

func (s *DeploymentStore) List(ctx context.Context, applicationID int64) ([]*types.Deployment, error) {
	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*deployment{}

	stmt := database.Builder.
		Select(deploymentColumns).
		From("deployments").
		Where("deployment_application_id = ?", applicationID).
		OrderBy("deployment_created DESC")

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query to sql")
	}

	if err := db.SelectContext(ctx, &dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "List deployments by application id query failed")

	}
	return s.mapDBDeployments(dst)

}

func (s *DeploymentStore) Update(ctx context.Context, deployment *types.Deployment) error {
	const deploymentUpdate = `UPDATE deployments 
	SET 
		deployment_status = :deployment_status
		,deployment_error = :deployment_error
		,deployment_version = :deployment_version
		,deployment_machine = :deployment_machine
		,deployment_started = :deployment_started
		,deployment_stopped = :deployment_stopped
	WHERE deployment_id = :deployment_id AND deployment_version = :deployment_version-1`

	updatedTime := time.Now().UTC().UnixMilli()

	deployment.Version++
	deployment.Updated = updatedTime

	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(deploymentUpdate, deployment)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind deployment object")
	}

	request, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to update deployment")
	}

	count, err := request.RowsAffected()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to get update count")
	}

	if count == 0 {
		return baseStore.ErrVersionConflict
	}

	return nil

}

func (s *DeploymentStore) ListIncomplete(ctx context.Context) ([]*types.Deployment, error) {
	const sqlQuery = deploymentBaseQuery + `
			WHERE deployment_status IN ('pending', 'running')
			ORDER BY deployment_id ASC`

	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*deployment{}

	if err := db.SelectContext(ctx, &dst, sqlQuery); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select incomplete deployments query failed")
	}

	return s.mapDBDeployments(dst)
}

func (s *DeploymentStore) ListIncompleteByApplicationID(ctx context.Context, applicationID int64) ([]*types.Deployment, error) {
	const sqlQuery = deploymentBaseQuery + `
         WHERE deployment_application_id = $1
         AND deployment_status IN ('pending', 'running')
         ORDER BY deployment_id ASC`

	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*deployment{}

	if err := db.SelectContext(ctx, &dst, sqlQuery, applicationID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select incomplete deployments query failed")
	}

	return s.mapDBDeployments(dst)
}

func (s *DeploymentStore) mapDBDeployment(dbDeployment *deployment) (*types.Deployment, error) {
	d := &dbDeployment.Deployment

	compSpec := new(types.ApplicationSpec)
	if err := json.Unmarshal([]byte(d.SpecJson), &compSpec); err != nil {
		return nil, err
	}
	d.Spec = compSpec
	return d, nil
}

func (s *DeploymentStore) mapDBDeployments(dbDeployments []*deployment) ([]*types.Deployment, error) {
	dep := make([]*types.Deployment, len(dbDeployments))
	for i := range dbDeployments {
		d, err := s.mapDBDeployment(dbDeployments[i])
		if err != nil {
			return nil, err
		}
		dep[i] = d
	}
	return dep, nil
}
