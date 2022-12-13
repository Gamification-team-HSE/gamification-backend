package stat

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/stat"
)

type Service interface {
	Get(ctx context.Context, id int64) (*models.Stat, error)
	Create(ctx context.Context, user *models.Stat) error
	List(ctx context.Context) ([]*models.Stat, error)
}
type service struct {
	validate *validator.Validate
	statRepo stat.Repository
}

func (s *service) List(ctx context.Context) ([]*models.Stat, error) {
	stats, err := s.statRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (s *service) Create(ctx context.Context, stat *models.Stat) error {
	err := s.validate.Struct(stat)
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

func (s *service) Get(ctx context.Context, id int64) (*models.Stat, error) {
	return s.statRepo.Get(ctx, id)
}
