package resolvers

import (
	"context"
	"gitlab.com/krespix/gamification-api/internal/models"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
	"gitlab.com/krespix/gamification-api/pkg/utils"
)

func (r *Resolver) GetStat(ctx context.Context, id int) (*apiModels.Stat, error) {
	stat, err := r.statService.Get(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return modelsStatToAPI(stat), nil
}

func (r *Resolver) GetStats(ctx context.Context) ([]*apiModels.Stat, error) {
	stats, err := r.statService.List(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*apiModels.Stat, 0, len(stats))
	for _, u := range stats {
		res = append(res, modelsStatToAPI(u))
	}
	return res, nil
}

func (r *Resolver) CreateStat(ctx context.Context, stat apiModels.NewStat) (interface{}, error) {
	mStat := &models.Stat{
		Name: stat.Name,
	}
	if stat.Name != "" {
		mStat.Name = stat.Name
	}

	err := r.statService.Create(ctx, mStat)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
	}, nil
}

func modelsStatToAPI(stat *models.Stat) *apiModels.Stat {
	return &apiModels.Stat{
		ID:          int(stat.ID),
		Name:        stat.Name,
		Description: utils.SqlNullStringToString(stat.Description),
		CreatedAt:   stat.CreatedAt,
		StartAt:     stat.StartAt,
		Period:      stat.Period,
		SeqPeriod:   utils.SqlNullStringToString(stat.SeqPeriod),
	}
}
