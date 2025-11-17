CREATE TABLE logs(
log_deployment_id 	INTEGER PRIMARY KEY
,log_data 				BLOB NOT NULL

,UNIQUE (log_deployment_id)
,constraint fk_log_deployment_id FOREIGN KEY (log_deployment_id)
   REFERENCES deployments(deployment_id) MATCH SIMPLE
   ON UPDATE NO ACTION
   ON DELETE CASCADE
);
