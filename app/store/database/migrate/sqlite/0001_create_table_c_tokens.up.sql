CREATE TABLE tokens (
 token_id             INTEGER PRIMARY KEY AUTOINCREMENT
,token_type           TEXT
,token_uid            TEXT
,token_principal_id   INTEGER
,token_expires_at     BIGINT
,token_issued_at      BIGINT
,token_created_by     INTEGER

,CONSTRAINT fk_token_principal_id FOREIGN KEY (token_principal_id)
    REFERENCES principals (principal_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
);

-- create explicit unique index with case insensitivity
CREATE UNIQUE INDEX tokens_principal_id_uid ON tokens(token_principal_id, LOWER(token_uid));

-- recreate old indices if needed (principal_id can be ignored since above index includes it)
CREATE INDEX tokens_type_expires_at ON tokens(token_type, token_expires_at);