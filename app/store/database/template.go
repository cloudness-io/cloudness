package database

import (
	"context"
	"encoding/json"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ store.TemplateStore = (*TemplateStore)(nil)

// NewTemplateStore returns new TemplateStore
func NewTemplateStore(db *sqlx.DB) *TemplateStore {
	return &TemplateStore{db: db}
}

// TemplateStore implements a TemplateStore backed by a relational database.
type TemplateStore struct {
	db *sqlx.DB
}

// template is a DB representation of a template
type template struct {
	types.Template
}

// templateColumns defines the columns of the template table
var templateColumns = `
	template_id
	,template_slug
	,template_name
	,template_icon
	,template_readme
	,template_tags
	,template_created
	`

// templateWithSpec defines the columns of the template table with spec
var templateColumnsWithSpec = templateColumns + `
		,template_spec`

// Find the template by id
func (s *TemplateStore) Find(ctx context.Context, id int64) (*types.Template, error) {
	stmt := database.Builder.Select(templateColumnsWithSpec).
		From("templates").
		Where("template_id = ?", id)

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to convert find template query to sql")
	}

	dst := new(template)
	db := dbtx.GetAccessor(ctx, s.db)
	if err := db.GetContext(ctx, dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "select query by id failed")
	}

	return s.mapDBTemplateSpec(dst)
}

// UpsertMany updates or inserts the templates
func (s *TemplateStore) UpsertMany(ctx context.Context, templates []*types.Template) error {
	stmt := database.Builder.Insert("templates").
		Columns(`template_slug
					,template_name
					,template_icon
					,template_readme
					,template_tags
					,template_spec
					,template_created`)
	for _, t := range templates {
		stmt = stmt.Values(
			t.Slug,
			t.Name,
			t.Icon,
			t.ReadMe,
			t.Tags,
			t.SpecJson,
			t.Created,
		)
	}

	stmt = stmt.Suffix(`ON CONFLICT (template_slug) 
	DO UPDATE SET 
		template_name = EXCLUDED.template_name
		,template_icon = EXCLUDED.template_icon
		,template_readme = EXCLUDED.template_readme
		,template_tags = EXCLUDED.template_tags
		,template_spec = EXCLUDED.template_spec`)

	query, args, err := stmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to convert upsert  template query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "template upsert query failed")
	}

	return nil
}

// List lists the templates
func (s *TemplateStore) List(ctx context.Context) ([]*types.Template, error) {
	stmt := database.Builder.Select(templateColumns).
		From("templates")

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to convert list  template query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*template{}

	if err := db.SelectContext(ctx, &dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "list template query failed")
	}
	return s.mapDBTemplates(dst), nil
}

func (s *TemplateStore) ListBySlugs(ctx context.Context, slugs []string) ([]*types.Template, error) {
	stmt := database.Builder.Select(templateColumns).
		From("templates").
		Where(sq.Eq{"template_slug": slugs})

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to convert list  template query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*template{}

	if err := db.SelectContext(ctx, &dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "list template query failed")
	}
	return s.mapDBTemplates(dst), nil
}

func (s *TemplateStore) mapDBTemplates(dst []*template) []*types.Template {
	var templates []*types.Template
	for _, t := range dst {
		templates = append(templates, &t.Template)
	}
	return templates
}

func (s *TemplateStore) mapDBTemplateSpec(dst *template) (*types.Template, error) {
	template := &dst.Template

	templateSpec := new(types.TemplateSpec)
	if err := json.Unmarshal([]byte(template.SpecJson), &templateSpec); err != nil {
		return nil, err
	}

	template.Spec = templateSpec
	return template, nil
}
