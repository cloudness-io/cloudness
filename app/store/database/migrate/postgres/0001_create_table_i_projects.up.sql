CREATE TABLE projects (
    project_id SERIAL PRIMARY KEY,
    project_uid INTEGER NOT NULL,
    project_tenant_id INTEGER REFERENCES tenants (tenant_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    project_name TEXT NOT NULL,
    project_description TEXT,
    project_created_by INTEGER NOT NULL,
    project_created BIGINT NOT NULL,
    project_updated BIGINT NOT NULL,
    project_deleted BIGINT DEFAULT NULL,
    UNIQUE (
        project_uid,
        project_tenant_id
    )
);