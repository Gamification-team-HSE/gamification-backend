package resolvers

import (
	"context"
	"database/sql"
	"time"

	"gitlab.com/krespix/gamification-api/internal/models"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
	"gitlab.com/krespix/gamification-api/pkg/utils"
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

func (r *Resolver) DeleteStat(ctx context.Context, id int) (interface{}, error) {
	err := r.statService.Delete(ctx, id)
	return nil, err
}

func (r *Resolver) GetStat(ctx context.Context, id int) (*apiModels.Stat, error) {
	stat, err := r.statService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &apiModels.Stat{
		ID:          int(stat.ID),
		Name:        stat.Name,
		Description: utils.SqlNullStringToString(stat.Description),
		CreatedAt:   int(stat.CreatedAt.Unix()),
		StartAt:     int(stat.StartAt.Unix()),
		Period:      stat.Period,
		SeqPeriod:   utils.SqlNullStringToString(stat.SeqPeriod),
	}, nil
}

func (r *Resolver) GetStats(ctx context.Context, pagination *apiModels.Pagination) (*apiModels.GetStatsResponse, error) {
	var (
		mPagination *models.Pagination
	)
	if pagination != nil {
		mPagination = &models.Pagination{
			Page: pagination.Page,
			Size: pagination.Size,
		}
	}
	res, err := r.statService.List(ctx, mPagination)
	if err != nil {
		return nil, err
	}
	resStats := make([]*apiModels.Stat, 0, len(res.Stats))
	for _, s := range res.Stats {
		resStats = append(resStats, &apiModels.Stat{
			ID:          int(s.ID),
			Name:        s.Name,
			Description: utils.SqlNullStringToString(s.Description),
			CreatedAt:   int(s.CreatedAt.Unix()),
			StartAt:     int(s.StartAt.Unix()),
			Period:      s.Period,
			SeqPeriod:   utils.SqlNullStringToString(s.SeqPeriod),
		})
	}

	return &apiModels.GetStatsResponse{
		Total: res.Total,
		Stats: resStats,
	}, nil
}

func (r *Resolver) UpdateStat(ctx context.Context, stat apiModels.UpdateStat) (interface{}, error) {
	sStat := &models.UpdateStat{
		ID: stat.ID,
	}
	if stat.Name != nil {
		sStat.Name = *stat.Name
	}
	if stat.Description != nil {
		sStat.Description = *stat.Description
	}
	if stat.StartAt != nil {
		sStat.StartedAt = time.Unix(int64(*stat.StartAt), 0)
	}
	if stat.Period != nil {
		sStat.Period = *stat.Period
	}
	if stat.SeqPeriod != nil {
		sStat.SeqPeriod = *stat.SeqPeriod
	}

	err := r.statService.Update(ctx, sStat)
	return nil, err
}
