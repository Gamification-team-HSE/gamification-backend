package image

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

type Service interface {
	Validate(ctx context.Context, image *graphql.Upload) error
}

type service struct {
}

func (s *service) Validate(ctx context.Context, image *graphql.Upload) error {
	//TODO доделать!!!!
	return nil
}

func New() Service {
	return &service{}
}
