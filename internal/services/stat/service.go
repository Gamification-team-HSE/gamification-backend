package stat

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/achievement"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/stat"
)

type Service interface {
	Create(ctx context.Context, user *models.Stat) error
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (*models.Stat, error)
	List(ctx context.Context, pagination *models.Pagination) (*models.GetStatsResponse, error)
	Update(ctx context.Context, updateStat *models.UpdateStat) error
}
type service struct {
	validate *validator.Validate
	statRepo stat.Repository
	achRepo  achievement.Repository
}

func (s *service) Update(ctx context.Context, updateStat *models.UpdateStat) error {
	return s.statRepo.Update(ctx, updateStat)
}

func (s *service) List(ctx context.Context, pagination *models.Pagination) (*models.GetStatsResponse, error) {
	total, err := s.statRepo.Total(ctx)
	if err != nil {
		return nil, err
	}
	var repoPagination *models.RepoPagination
	if pagination != nil {
		repoPagination = pagination.ToRepo()
	}
	stats, err := s.statRepo.List(ctx, repoPagination)
	if err != nil {
		return nil, err
	}
	return &models.GetStatsResponse{
		Stats: stats,
		Total: total,
	}, nil
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

func (s *service) Delete(ctx context.Context, id int) error {
	achList, err := s.achRepo.List(ctx, nil)
	if err != nil {
		return err
	}
	toUpdate := make([]*models.RepoAchievement, 0)

	needUpd := false
	for i, a := range achList {
		updRules := achList[i].Rules
		for j, b := range a.Rules.Blocks {
			res := make([]*models.StatRule, 0)
			for i, statRule := range b.StatRules {
				if statRule.StatID != id {
					res = append(res, b.StatRules[i])
				} else {
					needUpd = true
				}
			}
			updRules.Blocks[j].StatRules = res
		}
		if needUpd {
			toUpdate = append(toUpdate, &models.RepoAchievement{
				ID:    a.ID,
				Rules: updRules,
			})
		}
	}

	for _, a := range toUpdate {
		err = s.achRepo.Update(ctx, a)
		if err != nil {
			return err
		}
	}

	return s.statRepo.Delete(ctx, id)
}

func (s *service) Get(ctx context.Context, id int) (*models.Stat, error) {
	return s.statRepo.Get(ctx, id)
}

func New(
	statRepo stat.Repository,
	validate *validator.Validate,
	achRepo achievement.Repository,
) Service {
	return &service{
		statRepo: statRepo,
		validate: validate,
		achRepo:  achRepo,
	}
}
