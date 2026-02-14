CREATE TABLE applications (
    application_id SERIAL PRIMARY KEY,
    application_uid INTEGER NOT NULL,
    application_tenant_id INTEGER NOT NULL REFERENCES tenants (tenant_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    application_project_id INTEGER NOT NULL REFERENCES projects (project_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    application_environment_id INTEGER NOT NULL REFERENCES environments (environment_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    application_environment_uid INTEGER NOT NULL,
    application_server_id INTEGER NOT NULL REFERENCES servers (server_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    application_name TEXT NOT NULL,
	 application_slug TEXT NOT NULL,
	 application_parent_slug TEXT NOT NULL,
    application_description TEXT,
    application_domain TEXT,
    application_custom_domain TEXT,
    application_private_domain TEXT NOT NULL,
    application_status TEXT,
    application_type TEXT,
    application_spec TEXT NOT NULL,
    application_githubapp_id INTEGER DEFAULT NULL REFERENCES github_apps (github_app_id) ON UPDATE NO ACTION ON DELETE RESTRICT,
	 application_deployment_id INTEGER DEFAULT NULL,
    application_deployment_status TEXT,
    application_deployment_triggered_at BIGINT,
    application_created BIGINT NOT NULL,
    application_updated BIGINT NOT NULL,
    application_deleted BIGINT DEFAULT NULL,
    UNIQUE (
        application_tenant_id,
        application_project_id,
        application_environment_id,
		  application_slug
    ),
    UNIQUE (
        application_tenant_id,
        application_project_id,
        application_environment_id,
        application_private_domain
    )
);
