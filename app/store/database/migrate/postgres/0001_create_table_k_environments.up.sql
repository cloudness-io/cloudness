CREATE TABLE environments (
    environment_id SERIAL PRIMARY KEY,
    environment_uid INTEGER NOT NULL,
    environment_tenant_id INTEGER NOT NULL,
    environment_project_id INTEGER NOT NULL REFERENCES projects (project_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    environment_sequence INTEGER NOT NULL,
    environment_name TEXT NOT NULL,
	 environment_slug TEXT NOT NULL,
    environment_created_by INTEGER NOT NULL,
    environment_created BIGINT NOT NULL,
    environment_updated BIGINT NOT NULL,
    environment_deleted BIGINT DEFAULT NULL,
	 UNIQUE (
		  environment_tenant_id,
		  environment_project_id,
          environment_deleted,
		  environment_sequence		  
	 ),
	 UNIQUE (
		  environment_slug
	 )
);
