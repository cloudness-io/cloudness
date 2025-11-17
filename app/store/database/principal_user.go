package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	"github.com/pkg/errors"
)

// user is a DB representation of a user principal.
// It is required to allow storing transformed UIDs used for uniquness constraints and searching.
type user struct {
	types.User
}

const userColumns = principalCommonColumns + `
	,principal_user_password`

const userSelectBase = `
	SELECT` + userColumns + `
	FROM principals`

// FindUser finds the user by id.
func (s *PrincipalStore) FindUser(ctx context.Context, id int64) (*types.User, error) {
	const sqlQuery = userSelectBase + `
		WHERE principal_type = 'user' AND principal_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(user)
	if err := db.GetContext(ctx, dst, sqlQuery, id); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}

	return s.mapDBUser(dst), nil
}

// FindUserByUID finds the user by uid.
func (s *PrincipalStore) FindUserByUID(ctx context.Context, uid string) (*types.User, error) {
	const sqlQuery = userSelectBase + `
		WHERE principal_type = 'user' AND principal_uid = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(user)
	if err := db.GetContext(ctx, dst, sqlQuery, uid); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by uid query failed")
	}

	return s.mapDBUser(dst), nil
}

// FindUserByEmail finds the user by email.
func (s *PrincipalStore) FindUserByEmail(ctx context.Context, email string) (*types.User, error) {
	const sqlQuery = userSelectBase + `
		WHERE principal_type = 'user' AND LOWER(principal_email) = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	dst := new(user)
	if err := db.GetContext(ctx, dst, sqlQuery, strings.ToLower(email)); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by email query failed")
	}

	return s.mapDBUser(dst), nil
}

// CreateUser saves the user details.
func (s *PrincipalStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	const sqlQuery = `
		INSERT INTO principals (
			principal_type
			,principal_uid
			,principal_email
			,principal_display_name
			,principal_avatar_url
			,principal_blocked
			,principal_user_password
			,principal_salt
			,principal_created
			,principal_updated
		) values (
			'user'
			,:principal_uid
			,:principal_email
			,:principal_display_name
			,:principal_avatar_url
			,:principal_blocked
			,:principal_user_password
			,:principal_salt
			,:principal_created
			,:principal_updated
		) RETURNING principal_id`

	dbUser, err := s.mapToDBUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to map db user: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(sqlQuery, dbUser)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind user object")
	}

	if err = db.QueryRowContext(ctx, query, arg...).Scan(&user.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert query failed")
	}

	return user, nil
}

// UpdateUser updates an existing user.
func (s *PrincipalStore) UpdateUser(ctx context.Context, user *types.User) error {
	user.Updated = time.Now().UTC().UnixMilli()
	const sqlQuery = `
		UPDATE principals
		SET
			 principal_uid	          = :principal_uid
			,principal_email          = :principal_email
			,principal_display_name   = :principal_display_name
			,principal_avatar_url     = :principal_avatar_url
			,principal_blocked        = :principal_blocked
			,principal_updated        = :principal_updated
		WHERE principal_type = 'user' AND principal_id = :principal_id`

	dbUser, err := s.mapToDBUser(user)
	if err != nil {
		return fmt.Errorf("failed to map db user: %w", err)
	}

	db := dbtx.GetAccessor(ctx, s.db)

	query, arg, err := db.BindNamed(sqlQuery, dbUser)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Failed to bind user object")
	}

	if _, err = db.ExecContext(ctx, query, arg...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "Update query failed")
	}

	return err
}

// DeleteUser deletes the user.
func (s *PrincipalStore) DeleteUser(ctx context.Context, id int64) error {
	const sqlQuery = `
		DELETE FROM principals
		WHERE principal_type = 'user' AND principal_id = $1`

	db := dbtx.GetAccessor(ctx, s.db)

	if _, err := db.ExecContext(ctx, sqlQuery, id); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "The delete query failed")
	}

	return nil
}

// CountUsers counts the number of users.
func (s *PrincipalStore) CountUsers(ctx context.Context) (int64, error) {
	stmt := database.Builder.
		Select("count(1)").
		From("principals").
		Where("principal_type = 'user'").
		Where("principal_email != 'demo@cloudness.io'")

	sql, args, err := stmt.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	var count int64
	err = db.QueryRowContext(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "Failed executing count query")
	}
	return count, nil
}

func (s *PrincipalStore) mapDBUser(dbUser *user) *types.User {
	return &dbUser.User
}

func (s *PrincipalStore) mapDBUsers(dbUsers []*user) []*types.User {
	res := make([]*types.User, len(dbUsers))
	for i := range dbUsers {
		res[i] = s.mapDBUser(dbUsers[i])
	}
	return res
}

func (s *PrincipalStore) mapToDBUser(usr *types.User) (*user, error) {
	// user comes from outside.
	if usr == nil {
		return nil, fmt.Errorf("user is nil")
	}

	dbUser := &user{
		User: *usr,
	}

	return dbUser, nil
}
