package resolvers

import (
	"context"
	"gitlab.com/krespix/gamification-api/internal/models"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
)

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
