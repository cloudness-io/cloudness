CREATE TABLE environments (
 environment_id           INTEGER PRIMARY KEY AUTOINCREMENT
,environment_uid          INTEGER NOT NULL
,environment_tenant_id    INTEGeR NOT NULL
,environment_project_id   INTEGER NOT NULL
,environment_sequence     INTEGER
,environment_name         TEXT NOT NULL
,environment_slug         TEXT NOT NULL
,environment_created_by   INTEGER NOT NULL
,environment_created      BIGINT NOT NULL
,environment_updated      BIGINT NOT NULL
,environment_deleted      BIGINT DEFAULT NULL

,UNIQUE(environment_tenant_id, environment_project_id, environment_deleted, environment_sequence)
,UNIQUE(environment_slug)

,CONSTRAINT fk_environment_project_id FOREIGN KEY (environment_project_id)
    REFERENCES projects (project_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
);
