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

var _ store.FavoriteStore = (*FavoriteStore)(nil)

func NewFavoriteStore(db *sqlx.DB) *FavoriteStore {
	return &FavoriteStore{db}
}

type FavoriteStore struct {
	db *sqlx.DB
}

type favorite struct {
	types.Favorite
}

type favoriteDTO struct {
	ProjectUID      int64  `db:"project_uid"`
	ProjectName     string `db:"project_name"`
	EnvironmentUID  int64  `db:"environment_uid"`
	EnvironmentName string `db:"environment_name"`
	ApplicationUID  int64  `db:"application_uid"`
	AppName         string `db:"application_name"`
	AppDomain       string `db:"application_domain"`
}

// Add adds a favorite
func (s *FavoriteStore) Add(ctx context.Context, userID, applicationID int64) error {
	stmt := database.Builder.Insert("favorites").
		Columns(`favorite_user_id 
				,favorite_application_id
				,favorite_created`).
		Values(
			userID,
			applicationID,
			time.Now().UTC().UnixMilli(),
		)

	query, args, err := stmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to build favorite insert query")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to insert favorite")
	}

	return nil
}

// Delete delete a favorite
func (s *FavoriteStore) Delete(ctx context.Context, userID, applicationID int64) error {
	stmt := database.Builder.
		Delete("favorites").
		Where("favorite_user_id = ?", userID).
		Where("favorite_application_id = ?", applicationID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to convert delete favorite query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)

	_, err = db.ExecContext(ctx, sql, args...)
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "the delete query failed")
	}

	return nil
}

// List lists the favorite by tenant id
func (s *FavoriteStore) List(ctx context.Context, userID, tenantID int64) ([]*types.FavoriteDTO, error) {
	stmt := database.Builder.Select(`
					projects.project_uid
					,projects.project_name
					,environments.environment_uid
					,environments.environment_name
					,applications.application_uid
					,applications.application_name
					,applications.application_domain
				`).
		From("favorites").
		InnerJoin("applications ON applications.application_id = favorites.favorite_application_id").
		InnerJoin("environments ON environments.environment_id = applications.application_environment_id").
		InnerJoin("projects ON projects.project_id = applications.application_project_id").
		InnerJoin("project_memberships ON project_memberships.project_membership_project_id = projects.project_id").
		Where("favorite_user_id = ?", userID).
		Where("applications.application_tenant_id = ?", tenantID).
		Where("applications.application_deleted IS NULL")

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to convert query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*favoriteDTO{}
	if err := db.SelectContext(ctx, &dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select favorite by user id and tenant id query failed")
	}

	return s.mapDBFavoritesDTO(dst), nil
}

// Find find the favorite by user id and application id
func (s *FavoriteStore) Find(ctx context.Context, userID, applicationID int64) (*types.Favorite, error) {
	stmt := database.Builder.
		Select(`favorite_user_id
				,favorite_application_id
				,favorite_created`).
		From("favorites").
		Where("favorite_user_id = ?", userID).
		Where("favorite_application_id = ?", applicationID)

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to convert favorite query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(favorite)
	if err := db.GetContext(ctx, dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select favorite by user id and app id query failed")
	}
	return s.mapDBFavorite(dst), nil

}

func (s *FavoriteStore) mapDBFavorite(dbFav *favorite) *types.Favorite {
	return &dbFav.Favorite
}

func (s *FavoriteStore) mapDBFavoritesDTO(dbFavs []*favoriteDTO) []*types.FavoriteDTO {
	dst := make([]*types.FavoriteDTO, len(dbFavs))
	for i, dbFav := range dbFavs {
		dst[i] = &types.FavoriteDTO{
			ProjectUID:      dbFav.ProjectUID,
			ProjectName:     dbFav.ProjectName,
			EnvironmentUID:  dbFav.EnvironmentUID,
			EnvironmentName: dbFav.EnvironmentName,
			ApplicationUID:  dbFav.ApplicationUID,
			AppName:         dbFav.AppName,
			AppDomain:       dbFav.AppDomain,
		}
	}

	return dst
}
