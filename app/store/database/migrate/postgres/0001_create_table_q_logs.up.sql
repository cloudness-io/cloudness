CREATE TABLE logs (
    log_deployment_id INTEGER PRIMARY KEY REFERENCES deployments (deployment_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    log_data BYTEA NOT NULL,
    UNIQUE (log_deployment_id)
);