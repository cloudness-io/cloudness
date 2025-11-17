CREATE TABLE favorites (
 favorite_user_id							INTEGER NOT NULL
,favorite_application_id				INTEGER NOT NULL

,favorite_created						BIGINT NOT NULL

,UNIQUE (favorite_user_id, favorite_application_id)

,CONSTRAINT fk_favorite_application_id FOREIGN KEY (favorite_application_id)
    REFERENCES applications (application_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
);
