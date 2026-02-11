package database

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/helpers"
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

type templateRow struct {
	ID       int64  `db:"template_id"`
	Slug     string `db:"template_slug"`
	Name     string `db:"template_name"`
	Icon     string `db:"template_icon"`
	ReadMe   string `db:"template_readme"`
	SpecJSON string `db:"template_spec"`
	Created  int64  `db:"template_created"`
	TagsJSON string `db:"tags_json"`
}

// Find the template by id
func (s *TemplateStore) Find(ctx context.Context, id int64) (*types.Template, error) {
	stmt := s.templateSelect(true).
		Where("t.template_id = ?", id)

	return s.fetchTemplate(ctx, stmt, true)
}

// UpsertMany updates or inserts the templates
func (s *TemplateStore) UpsertMany(ctx context.Context, templates []*types.Template) error {
	now := time.Now().UTC().UnixMilli()
	runner := dbtx.New(s.db)

	return runner.WithTx(ctx, func(ctx context.Context) error {
		db := dbtx.GetAccessor(ctx, s.db)

		for _, tmpl := range templates {
			if tmpl.Created == 0 {
				tmpl.Created = now
			}

			if err := s.upsertTemplate(ctx, db, tmpl, now); err != nil {
				return err
			}
		}

		return nil
	})
}

// List lists the templates
func (s *TemplateStore) List(ctx context.Context) ([]*types.Template, error) {
	stmt := s.templateSelect(false)
	return s.fetchTemplates(ctx, stmt, false)
}

func (s *TemplateStore) ListTags(ctx context.Context) ([]*types.Tag, error) {
	stmt := database.Builder.Select(
		"tag_id",
		"tag_slug",
		"tag_name",
		"tag_created",
	).
		From("tags").
		OrderBy("tag_name ASC")

	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to convert list tags query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	rows := []*types.Tag{}

	if err := db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "list tags query failed")
	}

	return rows, nil
}

func (s *TemplateStore) ListNotInSlugs(ctx context.Context, slugs []string) ([]*types.Template, error) {
	stmt := s.templateSelect(false).
		Where(sq.NotEq{"t.template_slug": slugs})

	return s.fetchTemplates(ctx, stmt, false)
}

func (s *TemplateStore) ListBySlugs(ctx context.Context, slugs []string) ([]*types.Template, error) {
	stmt := s.templateSelect(false).
		Where(sq.Eq{"t.template_slug": slugs})

	return s.fetchTemplates(ctx, stmt, false)
}

func (s *TemplateStore) ListByTag(ctx context.Context, tag string) ([]*types.Template, error) {
	slug := helpers.Normalize(tag)
	if slug == "" {
		return []*types.Template{}, nil
	}

	stmt := s.templateSelect(false).
		Where("t.template_id IN (SELECT tt.template_id FROM template_tags tt JOIN tags tg ON tt.tag_id = tg.tag_id WHERE tg.tag_slug = ?)", slug)

	return s.fetchTemplates(ctx, stmt, false)
}

func (s *TemplateStore) templateSelect(withSpec bool) sq.SelectBuilder {
	driver := s.db.DriverName()

	columns := []string{
		"t.template_id",
		"t.template_slug",
		"t.template_name",
		"t.template_icon",
		"t.template_readme",
		"t.template_created",
		tagAggregationExpression(driver) + " AS tags_json",
	}

	groupBy := []string{
		"t.template_id",
		"t.template_slug",
		"t.template_name",
		"t.template_icon",
		"t.template_readme",
		"t.template_created",
	}

	if withSpec {
		columns = append(columns, "t.template_spec")
		groupBy = append(groupBy, "t.template_spec")
	}

	return database.Builder.Select(columns...).
		From("templates t").
		LeftJoin("template_tags tt ON tt.template_id = t.template_id").
		LeftJoin("tags tg ON tg.tag_id = tt.tag_id").
		GroupBy(groupBy...)
}

func (s *TemplateStore) fetchTemplate(ctx context.Context, stmt sq.SelectBuilder, withSpec bool) (*types.Template, error) {
	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to convert find template query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := new(templateRow)
	if err := db.GetContext(ctx, dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "select query by id failed")
	}

	return s.mapTemplateRow(dst, withSpec)
}

func (s *TemplateStore) fetchTemplates(ctx context.Context, stmt sq.SelectBuilder, withSpec bool) ([]*types.Template, error) {
	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "failed to convert list template query to sql")
	}

	db := dbtx.GetAccessor(ctx, s.db)
	dst := []*templateRow{}

	if err := db.SelectContext(ctx, &dst, query, args...); err != nil {
		return nil, database.ProcessSQLErrorf(ctx, err, "list template query failed")
	}

	return s.mapTemplateRows(dst, withSpec)
}

func (s *TemplateStore) mapTemplateRows(dst []*templateRow, withSpec bool) ([]*types.Template, error) {
	templates := make([]*types.Template, 0, len(dst))
	for _, t := range dst {
		tmpl, err := s.mapTemplateRow(t, withSpec)
		if err != nil {
			return nil, err
		}
		templates = append(templates, tmpl)
	}
	return templates, nil
}

func (s *TemplateStore) mapTemplateRow(dst *templateRow, withSpec bool) (*types.Template, error) {
	tmpl := &types.Template{
		ID:       dst.ID,
		Slug:     dst.Slug,
		Name:     dst.Name,
		Icon:     dst.Icon,
		ReadMe:   dst.ReadMe,
		SpecJson: dst.SpecJSON,
		Created:  dst.Created,
	}

	if dst.TagsJSON != "" {
		if err := json.Unmarshal([]byte(dst.TagsJSON), &tmpl.Tags); err != nil {
			return nil, fmt.Errorf("failed to decode template tags: %w", err)
		}
	}

	if len(tmpl.Tags) > 1 {
		sort.Strings(tmpl.Tags)
	}

	if withSpec && dst.SpecJSON != "" {
		templateSpec := new(types.TemplateSpec)
		if err := json.Unmarshal([]byte(tmpl.SpecJson), templateSpec); err != nil {
			return nil, err
		}
		templateSpec.Tags = tmpl.Tags
		tmpl.Spec = templateSpec
	}

	return tmpl, nil
}

func (s *TemplateStore) upsertTemplate(ctx context.Context, db dbtx.Accessor, tmpl *types.Template, tagCreated int64) error {
	stmt := database.Builder.Insert("templates").
		Columns(`template_slug
                ,template_name
                ,template_icon
                ,template_readme
                ,template_spec
                ,template_created`).
		Values(
			tmpl.Slug,
			tmpl.Name,
			tmpl.Icon,
			tmpl.ReadMe,
			tmpl.SpecJson,
			tmpl.Created,
		)

	stmt = stmt.Suffix(`ON CONFLICT (template_slug) 
    DO UPDATE SET 
        template_name = EXCLUDED.template_name
        ,template_icon = EXCLUDED.template_icon
        ,template_readme = EXCLUDED.template_readme
        ,template_spec = EXCLUDED.template_spec`)

	query, args, err := stmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to convert upsert  template query to sql")
	}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "template upsert query failed")
	}

	templateID, err := s.templateIDBySlug(ctx, db, tmpl.Slug)
	if err != nil {
		return err
	}

	return s.replaceTemplateTags(ctx, db, templateID, tmpl.Tags, tagCreated)
}

func (s *TemplateStore) templateIDBySlug(ctx context.Context, db dbtx.Accessor, slug string) (int64, error) {
	stmt := database.Builder.Select("template_id").
		From("templates").
		Where("template_slug = ?", slug)

	query, args, err := stmt.ToSql()
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "failed to convert find template id query to sql")
	}

	var templateID int64
	if err := db.GetContext(ctx, &templateID, query, args...); err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "failed to get template id")
	}

	return templateID, nil
}

func (s *TemplateStore) replaceTemplateTags(ctx context.Context, db dbtx.Accessor, templateID int64, tags []string, created int64) error {
	deleteStmt := database.Builder.Delete("template_tags").
		Where("template_id = ?", templateID)

	query, args, err := deleteStmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to convert delete template tags query to sql")
	}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to delete template tags")
	}

	if len(tags) == 0 {
		return nil
	}

	tagIDs, err := s.ensureTags(ctx, db, tags, created)
	if err != nil {
		return err
	}

	insertStmt := database.Builder.Insert("template_tags").
		Columns("template_id", "tag_id")

	for _, tagID := range tagIDs {
		insertStmt = insertStmt.Values(templateID, tagID)
	}

	insertStmt = insertStmt.Suffix("ON CONFLICT (template_id, tag_id) DO NOTHING")

	query, args, err = insertStmt.ToSql()
	if err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to convert insert template tags query to sql")
	}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return database.ProcessSQLErrorf(ctx, err, "failed to insert template tags")
	}

	return nil
}

func (s *TemplateStore) ensureTags(ctx context.Context, db dbtx.Accessor, tags []string, created int64) ([]int64, error) {
	seen := make(map[string]struct{})
	ids := make([]int64, 0, len(tags))

	for _, tag := range tags {
		clean := strings.TrimSpace(tag)
		if clean == "" {
			continue
		}

		slug := helpers.Normalize(clean)
		if slug == "" {
			continue
		}

		if _, ok := seen[slug]; ok {
			continue
		}
		seen[slug] = struct{}{}

		tagID, err := s.upsertTag(ctx, db, slug, clean, created)
		if err != nil {
			return nil, err
		}

		ids = append(ids, tagID)
	}

	return ids, nil
}

func (s *TemplateStore) upsertTag(ctx context.Context, db dbtx.Accessor, slug string, name string, created int64) (int64, error) {
	stmt := database.Builder.Insert("tags").
		Columns(`tag_slug
                ,tag_name
                ,tag_created`).
		Values(
			slug,
			name,
			created,
		)

	stmt = stmt.Suffix(`ON CONFLICT (tag_slug) 
    DO UPDATE SET 
        tag_name = EXCLUDED.tag_name`)

	query, args, err := stmt.ToSql()
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "failed to convert upsert tag query to sql")
	}

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "tag upsert query failed")
	}

	selectStmt := database.Builder.Select("tag_id").
		From("tags").
		Where("tag_slug = ?", slug)

	query, args, err = selectStmt.ToSql()
	if err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "failed to convert find tag id query to sql")
	}

	var tagID int64
	if err := db.GetContext(ctx, &tagID, query, args...); err != nil {
		return 0, database.ProcessSQLErrorf(ctx, err, "failed to get tag id")
	}

	return tagID, nil
}

func tagAggregationExpression(driver string) string {
	switch driver {
	case "sqlite3":
		return "COALESCE(json_group_array(tg.tag_slug), '[]')"
	default:
		return "COALESCE(json_agg(tg.tag_slug ORDER BY tg.tag_slug) FILTER (WHERE tg.tag_slug IS NOT NULL), '[]')"
	}
}
