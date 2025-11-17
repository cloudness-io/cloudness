CREATE TABLE tenant_memberships (
    tenant_membership_id SERIAL PRIMARY KEY,
    tenant_membership_tenant_id INTEGER NOT NULL REFERENCES tenants (tenant_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    tenant_membership_principal_id INTEGER NOT NULL REFERENCES principals (principal_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    tenant_membership_role TEXT,
    tenant_membership_created_by INTEGER NOT NULL REFERENCES principals (principal_id) ON UPDATE NO ACTION ON DELETE NO ACTION,
    tenant_membership_created BIGINT NOT NULL,
    tenant_membership_updated BIGINT NOT NULL,
    UNIQUE (
        tenant_membership_tenant_id,
        tenant_membership_principal_id
    )
);