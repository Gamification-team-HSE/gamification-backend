package resolvers

import (
	"gitlab.com/krespix/gamification-api/internal/models"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
)

func (r *Resolver) CreateStat(stat apiModels.NewStat) (interface{}, error) {
	mStat := &models.Stat{
		Name: stat.Name,
	}
	if stat.Name != "" {
		mStat.Name = stat.Name
	}

	/* Это тесты?
	err := r.
	if err != nil {
		return nil, err
	}*/

	return map[string]interface{}{
		"status": "success",
	}, nil
}

/*
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
*/
