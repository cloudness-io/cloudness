CREATE TABLE variables (
    variable_uid SERIAL PRIMARY KEY,
    variable_environment_id INTEGER NOT NULL REFERENCES environments (environment_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    variable_application_id INTEGER NOT NULL REFERENCES applications (application_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    variable_key TEXT NOT NULL,
    variable_value TEXT,
    variable_text_value TEXT,
    variable_type TEXT NOT NULL,
    variable_created BIGINT NOT NULL,
    variable_updated BIGINT NOT NULL,
    UNIQUE (
        variable_environment_id,
        variable_application_id,
        variable_key
    )
);