CREATE TABLE templates (
    template_id SERIAL PRIMARY KEY,
    template_slug TEXT NOT NULL,
    template_name TEXT NOT NULL,
    template_icon TEXT,
    template_readme TEXT,
    template_spec TEXT NOT NULL,
    template_created BIGINT NOT NULL,
    UNIQUE (template_slug)
);

CREATE TABLE tags (
    tag_id SERIAL PRIMARY KEY,
    tag_slug TEXT NOT NULL,
    tag_name TEXT NOT NULL,
    tag_created BIGINT NOT NULL,
    UNIQUE (tag_slug)
);

CREATE TABLE template_tags (
    template_id INTEGER NOT NULL REFERENCES templates (template_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags (tag_id) ON UPDATE NO ACTION ON DELETE CASCADE,
    PRIMARY KEY (template_id, tag_id)
);

CREATE INDEX idx_template_tags_tag_id ON template_tags(tag_id);
CREATE INDEX idx_template_tags_template_id ON template_tags(template_id);