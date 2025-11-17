package database

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ store.VariableStore = (*VariableStore)(nil)

// VariableStore implements a VariableStore backed by a relational database.
type VariableStore struct {
	db *sqlx.DB
}

// NewVariableStore creates a new VariableStore
func NewVariableStore(db *sqlx.DB) *VariableStore {
	return &VariableStore{
		db: db,
	}
}

// variable is a DB representation of a  variable
type variable struct {
	types.Variable
}

// variableColumns defines the columns for  variable
var variableColumns = `
	variable_uid
	,variable_environment_id
	,variable_application_id
	,variable_key
	,variable_value
	,variable_text_value
	,variable_type
	,variable_created
	,variable_updated`

// variableAppRefColumns defines the columns for  variable application reference
var variableAppRefColumns = variableColumns + `
		,application_name`

// Find the variable by application id and variable uid
func (s *VariableStore) Find(ctx context.Context, applicationID, varUID int64) (*types.Variable, error) {
	stmt := database.Builder.Select(variableColumns).
		From("variables").
		Where("variable_application_id = ? AND variable_uid = ?", applicationID, varUID)

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to convert  variable query to sql")
	}

	dst := new(variable)
	db := dbtx.GetAccessor(ctx, s.db)
	if err := db.GetContext(ctx, dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select query failed")
	}

	return s.mapDBVariable(dst), nil
}

// Upsert updates or inserts the  variables
func (s *VariableStore) Upsert(ctx context.Context, variable *types.Variable) error {
	return s.UpsertMany(ctx, []*types.Variable{variable})
}

// UpsertMany updates or inserts the  variables
func (s *VariableStore) UpsertMany(ctx context.Context, variables []*types.Variable) error {
	stmt := database.Builder.Insert("variables").
		Columns(variableColumns)
	for _, v := range variables {
		stmt = stmt.Values(
			v.UID,
			v.EnvironmentID,
			v.ApplicationID,
			v.Key,
			v.Value,
			v.TextValue,
			v.Type,
			v.Created,
			v.Updated,
		)
	}

	stmt = stmt.Suffix(`ON CONFLICT (variable_environment_id, variable_application_id, variable_key) 
	DO UPDATE SET 
	variable_key = EXCLUDED.variable_key
	,variable_value = EXCLUDED.variable_value
	,variable_text_value = EXCLUDED.variable_text_value
	,variable_type = EXCLUDED.variable_type
	,variable_updated = ?`, time.Now().UTC().UnixMilli())

	query, args, err := stmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to convert upsert  variable query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, " variable upsert query failed")
	}

	return nil
}

// List lists the  variables by environment id and application id
func (s *VariableStore) List(ctx context.Context, environmentID, applicaitonID int64) ([]*types.Variable, error) {
	stmt := database.Builder.Select(variableColumns).
		From("variables").
		Where("variable_environment_id = ?", environmentID).
		Where("variable_application_id = ?", applicaitonID)

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to bind  variables object for list")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*variable{}

	if err := db.SelectContext(ctx, &dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to list  variables")
	}

	return s.mapDBVariables(dst), nil
}

func (s *VariableStore) ListInEnvironment(ctx context.Context, envID int64) ([]*types.Variable, error) {
	stmt := database.Builder.Select(variableAppRefColumns).
		From("variables").
		LeftJoin("applications ON applications.application_id = variables.variable_application_id").
		Where("variable_environment_id = ?", envID).
		Where("application_deleted is NULL")

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to bind  variables object for list in environment")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*variable{}
	if err := db.SelectContext(ctx, &dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to list  variables in environment")
	}

	return s.mapDBVariables(dst), nil
}

func (s *VariableStore) Delete(ctx context.Context, applicationID, varUID int64) error {
	stmt := database.Builder.Delete("variables").
		Where("variable_application_id = ? AND variable_uid = ?", applicationID, varUID)

	query, args, err := stmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to bind  variables object for delete above")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to delete  variables above")
	}
	return nil
}

func (s *VariableStore) DeleteByKey(ctx context.Context, applicationID int64, key string) error {
	stmt := database.Builder.Delete("variables").
		Where("variable_application_id = ? AND variable_key = ?", applicationID, key)

	query, args, err := stmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to bind  variables object for delete by key")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to delete  variables by key")
	}
	return nil
}

func (s *VariableStore) DeleteByKeys(ctx context.Context, applicationID int64, keys []string) error {
	stmt := database.Builder.Delete("variables").
		Where("variable_application_id = ?", applicationID).
		Where(sq.Eq{"variable_key": keys})

	query, args, err := stmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to bind  variables object for delete by key")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to delete  variables by key")
	}
	return nil
}

func (s *VariableStore) mapDBVariables(dst []*variable) []*types.Variable {
	res := make([]*types.Variable, len(dst))
	for i := range dst {
		res[i] = s.mapDBVariable(dst[i])
	}
	return res
}

func (s *VariableStore) mapDBVariable(dst *variable) *types.Variable {
	return &dst.Variable
}
