CREATE TABLE auth_settings (
 auth_id                    INTEGER PRIMARY KEY AUTOINCREMENT
,auth_provider              TEXT NOT NULL
,auth_enabled               BOOLEAN
,auth_client_id             TEXT
,auth_client_secret         TEXT
,auth_base_url              TEXT
,auth_created               INTEGER
,auth_updated               INTEGER

,UNIQUE(auth_provider)
);
