CREATE TABLE tokens (
    token_id SERIAL PRIMARY KEY,
    token_type TEXT,
    token_uid TEXT,
    token_principal_id INTEGER REFERENCES principals (principal_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    token_expires_at BIGINT,
    token_issued_at BIGINT,
    token_created_by INTEGER
);

CREATE UNIQUE INDEX tokens_principal_id_uid ON tokens (
    token_principal_id,
    LOWER(token_uid)
);

CREATE INDEX tokens_type_expires_at ON tokens (token_type, token_expires_at);