CREATE TABLE applications (
 application_id                         INTEGER PRIMARY KEY AUTOINCREMENT
,application_uid                        INTEGER NOT NULL
,application_tenant_id                  INTEGER NOT NULL
,application_project_id                 INTEGER NOT NULL
,application_environment_id             INTEGER NOT NULL
,application_environment_uid            INTEGER NOT NULL
,application_server_id                  INTEGER NOT NULL
,application_name                       TEXT NOT NULL
,application_slug                       TEXT NOT NULL
,application_parent_slug                TEXT NOT NULL
,application_description                TEXT
,application_domain 					TEXT
,application_custom_domain				TEXT
,application_private_domain             TEXT NOT NULL
,application_status                     TEXT
,application_type                       TEXT
,application_spec                       TEXT NOT NULL
,application_githubapp_id               INTEGER DEFAULT NULL             -- Optional (nullable)
,application_deployment_id              INTEGER DEFAULT NULL             -- Optional (nullable)
,application_deployment_status          TEXT
,application_deployment_triggered_at    BIGINT

,application_created                    BIGINT NOT NULL
,application_updated                    BIGINT NOT NULL
,application_deleted                    BIGINT DEFAULT NULL

,UNIQUE (application_tenant_id, application_project_id, application_environment_id, application_slug)
,UNIQUE (application_tenant_id, application_project_id, application_environment_id, application_private_domain)

,CONSTRAINT fk_application_tenant_id FOREIGN KEY (application_tenant_id)
    REFERENCES tenants (tenant_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE


,CONSTRAINT fk_application_project_id FOREIGN KEY (application_project_id)
    REFERENCES projects (project_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE

,CONSTRAINT fk_application_environment_id FOREIGN KEY (application_environment_id)
    REFERENCES environments (environment_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE

,CONSTRAINT fk_application_server_id FOREIGN KEY (application_server_id)
    REFERENCES servers (server_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE

,CONSTRAINT fk_application_githubapp_id FOREIGN KEY (application_githubapp_id)
    REFERENCES github_apps (github_app_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE RESTRICT
);
