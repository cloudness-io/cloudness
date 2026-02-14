CREATE TABLE volumes (
 volume_id                  INTEGER PRIMARY KEY AUTOINCREMENT
,volume_uid                 INTEGER NOT NULL
,volume_tenant_id           INTEGER NOT NULL
,volume_project_id          INTEGER NOT NULL
,volume_environment_id      INTEGER NOT NULL
,volume_environment_uid     INTEGER NOT NULL
,volume_server_id           INTEGER NOT NULL
,volume_application_id      INTEGER DEFAULT NULL
,volume_name                TEXT NOT NULL
,volume_parent_slug         TEXT NOT NULL
,volume_slug                TEXT NOT NULL
,volume_mount_path          TEXT NOT NULL
,volume_host_path           TEXT
,volume_size                INTEGER NOT NULL
,volume_is_readonly         BOOLEAN

,volume_created             BIGINT NOT NULL
,volume_updated             BIGINT NOT NULL
,volume_deleted             BIGINT DEFAULT NULL

,UNIQUE(volume_uid)
,UNIQUE (volume_tenant_id, volume_project_id, volume_environment_id, volume_slug)

,CONSTRAINT fk_volume_server_id FOREIGN KEY (volume_server_id)
    REFERENCES servers (server_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE

,CONSTRAINT fk_volume_application_id FOREIGN KEY (volume_application_id)
    REFERENCES applications (application_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL
);
