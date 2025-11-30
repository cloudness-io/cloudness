CREATE TABLE instances (
 instance_id                            INTEGER PRIMARY KEY AUTOINCREMENT
,instance_super_admin						 INTEGER DEFAULT NULL
,instance_public_ipv4                   TEXT
,instance_public_ipv6                   TEXT
,instance_fqdn                          TEXT
,instance_update_enabled                BOOLEAN
,instance_update_check_frequency        TEXT
,instance_dns_validation_enabled        BOOLEAN
,instance_dns_servers                   TEXT
,instance_dns_provider						 TEXT
,instance_dns_provider_auth				 TEXT
,instance_user_signup_enabled 			 BOOLEAN
,instance_demo_user_enabled				 BOOLEAN
,instance_registry_enabled					 BOOLEAN
,instance_registry_size			 			 INTEGER 
,instance_registry_mirror_enabled 		 BOOLEAN
,instance_registry_mirror_size	 		 INTEGER
,instance_external_scripts					 TEXT
,instance_created                       INTEGER
,instance_updated                       INTEGER

,CONSTRAINT fk_instance_user_id FOREIGN KEY (instance_super_admin)
    REFERENCES principals (principal_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE RESTRICT
);
