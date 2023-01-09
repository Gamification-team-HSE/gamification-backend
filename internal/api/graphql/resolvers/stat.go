package resolvers

import (
	"context"
	"database/sql"
	"gitlab.com/krespix/gamification-api/internal/models"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
)

func (r *Resolver) CreateStat(ctx context.Context, stat apiModels.NewStat) (interface{}, error) {
	mStat := &models.Stat{
		Name: stat.Name,
	}
	if stat.Description != nil {
		mStat.Description = sql.NullString{
			String: *stat.Description,
			Valid:  true,
		}
	}
	if stat.SeqPeriod != nil {
		mStat.SeqPeriod = sql.NullString{
			String: *stat.SeqPeriod,
			Valid:  true,
		}
	}
	mStat.Period = stat.Period

	err := r.statService.Create(ctx, mStat)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
	}, nil
}
