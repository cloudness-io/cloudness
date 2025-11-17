CREATE TABLE github_apps (
 github_app_id              INTEGER PRIMARY KEY AUTOINCREMENT
,github_app_uid             INTEGER NOT NULL
,github_app_tenant_id       INTEGER
,github_app_project_id      INTEGER
,github_app_private_key_id  INTEGER
,github_app_is_tenant_wide  BOOLEAN
,github_app_name            TEXT NOT NULL
,github_app_organization    TEXT
,github_app_api_url         TEXT NOT NULL
,github_app_html_url        TEXT NOT NULL
,github_app_custom_user     TEXT NOT NULL
,github_app_custom_port     INTEGER NOT NULL
,github_app_app_id          INTEGER NULL
,github_app_installation_id INTEGER NULL
,github_app_client_id       TEXT NULL
,github_app_client_secret   TEXT NULL
,github_app_webhook_secret  TEXT NULL

,github_app_created_by   INTEGER NOT NULL
,github_app_created      BIGINT NOT NULL
,github_app_updated      BIGINT NOT NULL


,CONSTRAINT fk_github_app_tenant_id FOREIGN KEY (github_app_tenant_id)
    REFERENCES tenants (tenant_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
);
