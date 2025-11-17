CREATE TABLE tenant_memberships (
tenant_membership_id 					INTEGER PRIMARY KEY AUTOINCREMENT
,tenant_membership_tenant_id 			INTEGER NOT NULL
,tenant_membership_principal_id 		INTEGER NOT NULL
,tenant_membership_role 			   TEXT
,tenant_membership_created_by 		INTEGER NOT NULL
,tenant_membership_created        	BIGINT NOT NULL
,tenant_membership_updated        	BIGINT NOT NULL


,UNIQUE(tenant_membership_tenant_id, tenant_membership_principal_id)
,CONSTRAINT fk_tenant_membership_tenant_id FOREIGN KEY (tenant_membership_tenant_id)
    REFERENCES tenants (tenant_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_tenant_membership_principal_id FOREIGN KEY (tenant_membership_principal_id)
    REFERENCES principals (principal_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_tenant_membership_created_by FOREIGN KEY (tenant_membership_created_by)
    REFERENCES principals (principal_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
);
