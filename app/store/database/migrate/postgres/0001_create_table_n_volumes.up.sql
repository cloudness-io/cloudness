CREATE TABLE volumes (
    volume_id SERIAL PRIMARY KEY,
    volume_uid INTEGER NOT NULL,
    volume_tenant_id INTEGER NOT NULL REFERENCES tenants (tenant_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    volume_project_id INTEGER NOT NULL REFERENCES projects (project_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    volume_environment_id INTEGER NOT NULL REFERENCES environments (environment_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    volume_environment_uid INTEGER NOT NULL,
    volume_server_id INTEGER REFERENCES servers (server_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    volume_application_id INTEGER DEFAULT NULL REFERENCES applications (application_id) ON UPDATE NO ACTION ON DELETE SET NULL,
    volume_name TEXT NOT NULL,
    volume_mount_path TEXT NOT NULL,
    volume_host_path TEXT,
    volume_size INTEGER NOT NULL,
    volume_is_readonly BOOLEAN,
    volume_created BIGINT NOT NULL,
    volume_updated BIGINT NOT NULL,
    volume_deleted BIGINT DEFAULT NULL,
    UNIQUE (volume_uid)
);
