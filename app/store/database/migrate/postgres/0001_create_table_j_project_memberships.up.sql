CREATE TABLE project_memberships (
    project_membership_tenant_id INTEGER NOT NULL REFERENCES tenants (tenant_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    project_membership_tenant_membership_id INTEGER NOT NULL REFERENCES tenant_memberships (tenant_membership_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    project_membership_project_id INTEGER NOT NULL REFERENCES projects (project_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    project_membership_principal_id INTEGER NOT NULL REFERENCES principals (principal_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    project_membership_role TEXT NOT NULL,
    project_membership_created_by INTEGER NOT NULL REFERENCES principals (principal_id) ON UPDATE NO ACTION ON DELETE NO ACTION,
    project_membership_created BIGINT NOT NULL,
    project_membership_updated BIGINT NOT NULL,
    CONSTRAINT pk_memberships PRIMARY KEY (
        project_membership_tenant_id,
        project_membership_project_id,
        project_membership_principal_id
    )
);