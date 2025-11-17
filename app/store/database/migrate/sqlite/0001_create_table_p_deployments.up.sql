CREATE TABLE deployments (
 deployment_id                INTEGER PRIMARY KEY AUTOINCREMENT
,deployment_uid               INTEGER NOT NULL
,deployment_application_id 		INTEGER NOT NULL
,deployment_spec              TEXT NOT NULL
,deployment_needs_build 		BOOLEAN NOT NULL
,deployment_triggerer 			TEXT
,deployment_title 				TEXT
,deployment_action 				TEXT
,deployment_status				TEXT NOT NULL
,deployment_error					TEXT NOT NULL
,deployment_version				INTEGER NOT NULL
,deployment_machine 				TEXT 
,deployment_started				BIGINT NOT NULL
,deployment_stopped				BIGINT NOT NULL
,deployment_created           BIGINT NOT NULL
,deployment_updated           BIGINT NOT NULL


,UNIQUE (deployment_application_id, deployment_uid)
,CONSTRAINT fk_deployment_application_id FOREIGN KEY (deployment_application_id)
    REFERENCES applications (application_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
);

CREATE INDEX deployment_uid_application_id
ON deployments(deployment_uid, deployment_application_id, deployment_status);
