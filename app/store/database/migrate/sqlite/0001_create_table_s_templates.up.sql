--template--
CREATE TABLE templates (
 template_id					INTEGER PRIMARY KEY AUTOINCREMENT
,template_slug					TEXT NOT NULL
,template_name					TEXT NOT NULL
,template_icon					TEXT
,template_readme				TEXT
,template_spec					TEXT NOT NULL
,template_created				BIGINT NOT NULL

,UNIQUE(template_slug)
);


--tags--
CREATE TABLE tags (
 tag_id					INTEGER PRIMARY KEY AUTOINCREMENT		
,tag_slug				TEXT NOT NULL --lowercase name for unique detection
,tag_name				TEXT NOT NULL
,tag_created			BIGINT NOT NULL

,UNIQUE(tag_slug)
);

--tempalate <==> tag--
CREATE TABLE template_tags(
 template_id 			INTEGER NOT NULL
,tag_id                 INTEGER NOT NULL

,PRIMARY KEY (template_id, tag_id)
,CONSTRAINT fk_template_id FOREIGN KEY (template_id)
    REFERENCES templates (template_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
,CONSTRAINT fk_tag_id FOREIGN KEY (tag_id)
    REFERENCES tags (tag_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
);

CREATE INDEX idx_template_tags_tag_id ON template_tags(tag_id);
CREATE INDEX idx_template_tags_template_id ON template_tags(template_id);
