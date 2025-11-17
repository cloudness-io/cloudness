CREATE TABLE project_memberships (
 project_membership_tenant_id 					INTEGER NOT NULL
,project_membership_tenant_membership_id 		INTEGER NOT NULL
,project_membership_project_id 					INTEGER NOT NULL
,project_membership_principal_id 				INTEGER NOT NULL
,project_membership_role 							TEXT NOT NULL
,project_membership_created_by 					INTEGER NOT NULL
,project_membership_created 						BIGINT NOT NULL
,project_membership_updated 						BIGINT NOT NULL


,CONSTRAINT pk_memberships PRIMARY KEY (project_membership_tenant_id, project_membership_project_id, project_membership_principal_id)
,CONSTRAINT fk_project_membership_tenant_id FOREIGN KEY (project_membership_tenant_id)
    REFERENCES tenants (tenant_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_project_membership_tenant_membership_id FOREIGN KEY (project_membership_tenant_membership_id)
    REFERENCES tenant_memberships (tenant_membership_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_project_membership_project_id FOREIGN KEY (project_membership_project_id)
    REFERENCES projects (project_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_project_membership_principal_id FOREIGN KEY (project_membership_principal_id)
    REFERENCES principals (principal_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_project_membership_created_by FOREIGN KEY (project_membership_created_by)
    REFERENCES principals (principal_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
);
