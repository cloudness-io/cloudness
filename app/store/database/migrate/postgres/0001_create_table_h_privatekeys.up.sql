CREATE TABLE private_keys (
    private_key_id SERIAL PRIMARY KEY,
    private_key_uid INTEGER NOT NULL,
    private_key_tenant_id INTEGER REFERENCES tenants (tenant_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    private_key_name TEXT NOT NULL,
    private_key_description TEXT,
    private_key_pem TEXT,
    private_key_is_git BOOLEAN,
    private_key_created BIGINT NOT NULL,
    private_key_updated BIGINT NOT NULL,
    UNIQUE (
        private_key_uid,
        private_key_tenant_id
    )
);