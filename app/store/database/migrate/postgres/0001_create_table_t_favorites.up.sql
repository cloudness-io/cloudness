CREATE TABLE favorites (
    favorite_user_id INTEGER NOT NULL,
    favorite_application_id INTEGER NOT NULL REFERENCES applications (application_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    favorite_created BIGINT NOT NULL,
    UNIQUE (
        favorite_user_id,
        favorite_application_id
    )
);