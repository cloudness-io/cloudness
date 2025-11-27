package database

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	baseStore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ store.ServerStore = (*ServerStore)(nil)

// NewServerStore create a new server store.
func NewServerStore(db *sqlx.DB) store.ServerStore {
	return &ServerStore{db: db}
}

// ServerStore implements a ServerStore backed by a relational database.
type ServerStore struct {
	db *sqlx.DB
}

// server is a DB representation of a server
type server struct {
	types.Server
}

const serverColumns = `
	server_id 
	,server_uid 
	,server_type
	,server_name
	,server_description
	,server_ipv4
	,server_ipv6
	,server_wildcard_domain
	,server_dns_provider
	,server_dns_provider_auth
	,server_user 
	,server_port
	,server_volume_supports_online_expansion
	,server_builder_is_enabled
	,server_builder_is_build_server
	,server_builder_max_concurrent_builds
	,server_builder_max_cpu
	,server_builder_max_memory
	,server_volume_min_size
	,server_created
	,server_updated
	`

const serverInsert = `
INSERT INTO servers (
	server_uid
	,server_type
	,server_name
	,server_description
	,server_ipv4
	,server_ipv6
	,server_wildcard_domain
	,server_dns_provider
	,server_dns_provider_auth
	,server_user
	,server_port
	,server_volume_supports_online_expansion
	,server_builder_is_enabled
	,server_builder_is_build_server
	,server_builder_max_concurrent_builds
	,server_builder_max_cpu
	,server_builder_max_memory
	,server_volume_min_size
	,server_created
	,server_updated
) VALUES (
	:server_uid
	,:server_type
	,:server_name
	,:server_description
	,:server_ipv4
	,:server_ipv6
	,:server_wildcard_domain
	,:server_dns_provider
	,:server_dns_provider_auth
	,:server_user
	,:server_port
	,:server_volume_supports_online_expansion
	,:server_builder_is_enabled
	,:server_builder_is_build_server
	,:server_builder_max_concurrent_builds
	,:server_builder_max_cpu
	,:server_builder_max_memory
	,:server_volume_min_size
	,:server_created
	,:server_updated
) RETURNING server_id`

const serverUpdateStmt = `UPDATE servers 
	SET 
		server_name = :server_name
		,server_description = :server_description
		,server_ipv4 = :server_ipv4
		,server_ipv6 = :server_ipv6
		,server_wildcard_domain = :server_wildcard_domain
		,server_dns_provider = :server_dns_provider
		,server_dns_provider_auth = :server_dns_provider_auth
		,server_user = :server_user
		,server_port = :server_port
		,server_volume_supports_online_expansion = :server_volume_supports_online_expansion
		,server_builder_is_enabled = :server_builder_is_enabled
		,server_builder_is_build_server = :server_builder_is_build_server
		,server_builder_max_concurrent_builds = :server_builder_max_concurrent_builds
		,server_builder_max_cpu = :server_builder_max_cpu
		,server_builder_max_memory = :server_builder_max_memory
		,server_volume_min_size = :server_volume_min_size
		,server_updated = :server_updated
	WHERE server_id = :server_id
	`

// Find the server by id.
func (s *ServerStore) Find(ctx context.Context, id int64) (*types.Server, error) {
	stmt := database.Builder.
		Select(serverColumns).
		From("servers").
		Where("server_id = ?", id)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(server)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}
	return s.mapDBServer(dst), nil
}

// FindByUID the server by id.
func (s *ServerStore) FindByUID(ctx context.Context, serverUID int64) (*types.Server, error) {
	stmt := database.Builder.
		Select(serverColumns).
		From("servers").
		Where("server_uid = ?", serverUID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(server)

	if err := db.GetContext(ctx, dst, sql, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select by id query failed")
	}
	return s.mapDBServer(dst), nil
}

// Update  updates the server details
func (s *ServerStore) Update(ctx context.Context, server *types.Server) (*types.Server, error) {
	server.Updated = time.Now().UTC().UnixMilli()

	db := dbtx.GetAccessor(ctx, s.db)
	query, args, err := db.BindNamed(serverUpdateStmt, server)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind server object")
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Update query failed")
	}

	count, err := result.RowsAffected()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to get update count")
	}

	if count == 0 {
		return nil, baseStore.ErrVersionConflict
	}

	return server, nil
}

// Create save the server details
func (s *ServerStore) Create(ctx context.Context, server *types.Server) (*types.Server, error) {
	server.Created = time.Now().UTC().UnixMilli()
	db := dbtx.GetAccessor(ctx, s.db)

	query, args, err := db.BindNamed(serverInsert, server)
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Failed to bind server oject")
	}

	if err = db.QueryRowContext(ctx, query, args...).Scan(&server.ID); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Insert server query failed")
	}

	return server, nil
}

// List lists the server
func (s *ServerStore) List(ctx context.Context) ([]*types.Server, error) {
	stmt := database.Builder.
		Select(serverColumns).
		From("servers").
		OrderBy("server_created")

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dbServers := []*server{}
	if err := db.SelectContext(ctx, &dbServers, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "Select all query failed")
	}

	return s.mapDBServers(dbServers), nil
}

func (s *ServerStore) mapDBServer(d *server) *types.Server {
	return &d.Server
}

func (s *ServerStore) mapDBServers(dbServers []*server) []*types.Server {
	dep := make([]*types.Server, len(dbServers))
	for i := range dbServers {
		dep[i] = s.mapDBServer(dbServers[i])
	}
	return dep
}
