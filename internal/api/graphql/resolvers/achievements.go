package resolvers

import (
	"context"

	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
)

func (r *Resolver) DeleteAchievement(ctx context.Context, id int) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) CreateAchievement(ctx context.Context, achievement apiModels.CreateAchievement) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) UpdateAchievement(ctx context.Context, achievement apiModels.UpdateAchievement) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) GetAchievement(ctx context.Context, id int) (*apiModels.Achievement, error) {
	return nil, nil
}

func (r *Resolver) GetAchievements(ctx context.Context, pagination *apiModels.Pagination) (*apiModels.GetAchievementsResponse, error) {
	return nil, nil
}
