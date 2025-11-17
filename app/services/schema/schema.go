package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloudness-io/cloudness/schema"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"

	"github.com/qri-io/jsonschema"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service struct {
	applicationSchema *jsonschema.Schema
	templateSchema    *jsonschema.Schema
}

func NewService() *Service {
	s := &Service{
		applicationSchema: new(jsonschema.Schema),
		templateSchema:    new(jsonschema.Schema),
	}
	//application schema to qri-io schema
	if err := json.Unmarshal(schema.Application, s.applicationSchema); err != nil {
		log.Fatal().Err(err).Msg("Error parsing application schema")
	}

	if err := json.Unmarshal(schema.Template, s.templateSchema); err != nil {
		log.Fatal().Err(err).Msg("Error parsing template schema")
	}

	return s
}

func (s *Service) ValidateApplication(ctx context.Context, in *types.ApplicationSpec) error {
	inBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return s.validateBytes(ctx, s.applicationSchema, inBytes)
}

func (s *Service) ValidateTemplate(ctx context.Context, in []byte) error {
	return s.validateBytes(ctx, s.templateSchema, in)
}

func (s *Service) validateBytes(ctx context.Context, schema *jsonschema.Schema, in []byte) error {
	schemaErrors, err := schema.ValidateBytes(ctx, in)
	if err != nil {
		return err
	}

	if len(schemaErrors) > 0 {
		log.Ctx(ctx).Error().Any("Validation Error", schemaErrors).Msg("Validation Error")

		vErrs := check.NewValidationErrors()
		for _, sErr := range schemaErrors {
			key := getPropertyFromPath(sErr.PropertyPath)
			msg := s.sanitizeValidationError(key, sErr.Message)
			vErrs.AddValidationError(key, errors.New(msg))
		}
		return vErrs
	}

	return nil

}

func (s *Service) sanitizeValidationError(key string, msg string) string {
	switch msg {
	case `min length of 1 characters required: `:
		return fmt.Sprintf("%s is required", cases.Title(language.English, cases.NoLower).String(key))
	default:
		return msg
	}
}
