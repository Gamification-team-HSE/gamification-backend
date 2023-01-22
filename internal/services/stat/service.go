package stat

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/stat"
)

type Service interface {
	Create(ctx context.Context, user *models.Stat) error
}
type service struct {
	validate *validator.Validate
	statRepo stat.Repository
}

func (s *service) Create(ctx context.Context, stat *models.Stat) error {
	err := s.validate.Var(stat.Name, "required")
	if err != nil {
		return err
	}
	exists, err := s.statRepo.ExistsByName(ctx, stat.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("stat with name %s already exists", stat.Name)
	}
	return s.statRepo.Create(ctx, stat)
}

func New(
	statRepo stat.Repository,
	validate *validator.Validate,
) Service {
	return &service{
		statRepo: statRepo,
		validate: validate,
	}
}
