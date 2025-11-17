CREATE TABLE principals (
    principal_id SERIAL PRIMARY KEY,
    principal_uid TEXT,
    principal_email TEXT,
    principal_display_name TEXT,
    principal_avatar_url TEXT,
    principal_type TEXT,
    principal_blocked BOOLEAN,
    principal_user_password TEXT,
    principal_salt TEXT,
    principal_created BIGINT,
    principal_updated BIGINT,
    UNIQUE (principal_uid)
);

CREATE UNIQUE INDEX principals_lower_email ON principals (LOWER(principal_email));