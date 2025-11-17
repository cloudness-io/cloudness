CREATE TABLE projects (
 project_id             INTEGER PRIMARY KEY AUTOINCREMENT
,project_uid            INTEGER NOT NULL
,project_tenant_id      INTEGER
,project_name           TEXT NOT NULL
,project_description    TEXT
,project_created_by     INTEGER NOT NULL
,project_created        BIGINT NOT NULL
,project_updated        BIGINT NOT NULL
,project_deleted        BIGINT DEFAULT NULL

,UNIQUE(project_uid, project_tenant_id)

,CONSTRAINT fk_project_tenant_id FOREIGN KEY (project_tenant_id)
    REFERENCES tenants (tenant_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
);
