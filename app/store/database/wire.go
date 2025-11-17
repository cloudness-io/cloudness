package database

import (
	"context"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/app/store/database/migrate"
	"github.com/cloudness-io/cloudness/job"
	"github.com/cloudness-io/cloudness/store/database"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideDatabase,
	ProvideInstanceStore,
	ProvideServerStore,
	ProvideAuthSettingStore,
	ProvideUserStore,
	ProvideTokenStore,
	ProvideApplicationStore,
	ProvideTenantStore,
	ProvideTenantMembershipStore,
	ProvideProjectStore,
	ProvideProjectMembershipStore,
	ProvideEnvironmentStore,
	ProvideDeploymentStore,
	ProvideLogStore,
	ProvideGithubAppStore,
	ProvidePrivateKeyStore,
	ProvideVolumeStore,
	ProvideVariableStore,
	ProvideJobStore,
	ProvideTemplateStore,
	ProvideFavoriteStore,
)

// migrator is helper function to set up the database by performing automated
// database migration steps.
func migrator(ctx context.Context, db *sqlx.DB) error {
	return migrate.Migrate(ctx, db)
}

// ProvideDatabase provides a database connection.
func ProvideDatabase(ctx context.Context, config database.Config) (*sqlx.DB, error) {
	return database.ConnectAndMigrate(
		ctx,
		config.Driver,
		config.Datasource,
		migrator,
	)
}

// ProvideInstanceStore provides a instance store.
func ProvideInstanceStore(db *sqlx.DB) store.InstanceStore {
	return NewInstanceStore(db)
}

// ProvideServerStore provides a server store.
func ProvideServerStore(db *sqlx.DB) store.ServerStore {
	return NewServerStore(db)
}

// ProvideAuthSettingStore provides a auth setting store.
func ProvideAuthSettingStore(db *sqlx.DB) store.AuthSettingsStore {
	return NewAuthSettingsStore(db)
}

// ProvideUserStore provides a user store.
func ProvideUserStore(db *sqlx.DB) store.PrincipalStore {
	return NewPrincipalStore(db)
}

// ProvideTokenStore provides a token store.
func ProvideTokenStore(db *sqlx.DB) store.TokenStore {
	return NewTokenStore(db)
}

// ProvideApplicationStore provides a token store.
func ProvideApplicationStore(db *sqlx.DB) store.ApplicationStore {
	return NewApplicationStore(db)
}

// ProvideTenantStore provides a tenant store.
func ProvideTenantStore(db *sqlx.DB) store.TenantStore {
	return NewTenantStore(db)
}

// ProvideTenantMembershipRow provides a tenant membership store.
func ProvideTenantMembershipStore(db *sqlx.DB) store.TenantMembershipStore {
	return NewTenantMembershipStore(db)
}

// ProvideProjectStore provides a project store.
func ProvideProjectStore(db *sqlx.DB) store.ProjectStore {
	return NewProjectSore(db)
}

// ProvideProjectMembership provides a tenant membership store.
func ProvideProjectMembershipStore(db *sqlx.DB) store.ProjectMembershipStore {
	return NewProjectMembershipStore(db)
}

// ProvideEnvironmentStore provides a project store.
func ProvideEnvironmentStore(db *sqlx.DB) store.EnvironmentStore {
	return NewEnvironmentSore(db)
}

// ProvideDeploymentStore provides a deployment store.
func ProvideDeploymentStore(db *sqlx.DB) store.DeploymentStore {
	return NewDeploymentStore(db)
}

// ProvideLogStore provides a log store.
func ProvideLogStore(db *sqlx.DB) store.LogStore {
	return NewLogStore(db)
}

// ProvideGithubAppStore  provides a github app store.
func ProvideGithubAppStore(db *sqlx.DB) store.GithubAppStore {
	return NewGithubAppStore(db)
}

// ProvidePrivateKeyStore provides a private key store.
func ProvidePrivateKeyStore(db *sqlx.DB) store.PrivateKeyStore {
	return NewPrivateKeyStore(db)
}

// ProvideVolumeStore provides a volume store.
func ProvideVolumeStore(db *sqlx.DB) store.VolumeStore {
	return NewVolumeStore(db)
}

// ProvideEnvironmentVariableStore provides a environment variable store.
func ProvideVariableStore(db *sqlx.DB) store.VariableStore {
	return NewVariableStore(db)
}

// ProvideJobStore provides a volume store.
func ProvideJobStore(db *sqlx.DB) job.Store {
	return NewJobStore(db)
}

// ProvideTemplateStore provides a template store.
func ProvideTemplateStore(db *sqlx.DB) store.TemplateStore {
	return NewTemplateStore(db)
}

// ProvideFavoriteStore provides a favorite store.
func ProvideFavoriteStore(db *sqlx.DB) store.FavoriteStore {
	return NewFavoriteStore(db)
}
