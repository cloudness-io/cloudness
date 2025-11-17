CREATE TABLE auth_settings (
    auth_id SERIAL PRIMARY KEY,
    auth_provider TEXT NOT NULL,
    auth_enabled BOOLEAN,
    auth_client_id TEXT,
    auth_client_secret TEXT,
    auth_base_url TEXT,
    auth_created BIGINT,
    auth_updated BIGINT,
    UNIQUE (auth_provider)
);