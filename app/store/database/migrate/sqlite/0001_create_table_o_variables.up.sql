CREATE TABLE variables (
 variable_uid						INTEGER PRIMARY KEY
,variable_environment_id		INTEGER NOT NULL
,variable_application_id		INTEGER NOT NULL
,variable_key						TEXT NOT NULL
,variable_value					TEXT
,variable_text_value				TEXT
,variable_type             	TEXT NOT NULL
,variable_created					BIGINT NOT NULL
,variable_updated					BIGINT NOT NULL


,UNIQUE (variable_environment_id, variable_application_id, variable_key)
,CONSTRAINT fk_variable_application_id FOREIGN KEY (variable_application_id)
    REFERENCES applications (application_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_variable_environment_id FOREIGN KEY (variable_environment_id)
    REFERENCES environments (environment_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
);
