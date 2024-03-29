package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.26

import (
	"context"
	"fmt"

	"gitlab.com/krespix/gamification-api/pkg/graphql/models"
	"gitlab.com/krespix/gamification-api/pkg/graphql/server"
)

// SendCode is the resolver for the SendCode field.
func (r *mutationResolver) SendCode(ctx context.Context, email string) (interface{}, error) {
	panic(fmt.Errorf("not implemented: SendCode - SendCode"))
}

// VerifyCode is the resolver for the VerifyCode field.
func (r *mutationResolver) VerifyCode(ctx context.Context, email string, code int) (string, error) {
	panic(fmt.Errorf("not implemented: VerifyCode - VerifyCode"))
}

// CreateUser is the resolver for the CreateUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, user models.NewUser) (interface{}, error) {
	panic(fmt.Errorf("not implemented: CreateUser - CreateUser"))
}

// DeleteUser is the resolver for the DeleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id int) (interface{}, error) {
	panic(fmt.Errorf("not implemented: DeleteUser - DeleteUser"))
}

// BanUser is the resolver for the BanUser field.
func (r *mutationResolver) BanUser(ctx context.Context, id int) (interface{}, error) {
	panic(fmt.Errorf("not implemented: BanUser - BanUser"))
}

// RecoverUser is the resolver for the RecoverUser field.
func (r *mutationResolver) RecoverUser(ctx context.Context, id int) (interface{}, error) {
	panic(fmt.Errorf("not implemented: RecoverUser - RecoverUser"))
}

// UpdateUser is the resolver for the UpdateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, user *models.UpdateUser) (interface{}, error) {
	panic(fmt.Errorf("not implemented: UpdateUser - UpdateUser"))
}

// CreateStat is the resolver for the CreateStat field.
func (r *mutationResolver) CreateStat(ctx context.Context, stat models.NewStat) (interface{}, error) {
	panic(fmt.Errorf("not implemented: CreateStat - CreateStat"))
}

// UpdateStat is the resolver for the UpdateStat field.
func (r *mutationResolver) UpdateStat(ctx context.Context, stat models.UpdateStat) (interface{}, error) {
	panic(fmt.Errorf("not implemented: UpdateStat - UpdateStat"))
}

// DeleteStat is the resolver for the DeleteStat field.
func (r *mutationResolver) DeleteStat(ctx context.Context, id int) (interface{}, error) {
	panic(fmt.Errorf("not implemented: DeleteStat - DeleteStat"))
}

// UpdateEvent is the resolver for the UpdateEvent field.
func (r *mutationResolver) UpdateEvent(ctx context.Context, event models.UpdateEvent) (interface{}, error) {
	panic(fmt.Errorf("not implemented: UpdateEvent - UpdateEvent"))
}

// CreateEvent is the resolver for the CreateEvent field.
func (r *mutationResolver) CreateEvent(ctx context.Context, event models.NewEvent) (interface{}, error) {
	panic(fmt.Errorf("not implemented: CreateEvent - CreateEvent"))
}

// DeleteEvent is the resolver for the DeleteEvent field.
func (r *mutationResolver) DeleteEvent(ctx context.Context, id int) (interface{}, error) {
	panic(fmt.Errorf("not implemented: DeleteEvent - DeleteEvent"))
}

// DeleteAchievement is the resolver for the DeleteAchievement field.
func (r *mutationResolver) DeleteAchievement(ctx context.Context, id int) (interface{}, error) {
	panic(fmt.Errorf("not implemented: DeleteAchievement - DeleteAchievement"))
}

// CreateAchievement is the resolver for the CreateAchievement field.
func (r *mutationResolver) CreateAchievement(ctx context.Context, achievement models.CreateAchievement) (interface{}, error) {
	panic(fmt.Errorf("not implemented: CreateAchievement - CreateAchievement"))
}

// UpdateAchievement is the resolver for the UpdateAchievement field.
func (r *mutationResolver) UpdateAchievement(ctx context.Context, achievement models.UpdateAchievement) (interface{}, error) {
	panic(fmt.Errorf("not implemented: UpdateAchievement - UpdateAchievement"))
}

// GetUser is the resolver for the GetUser field.
func (r *queryResolver) GetUser(ctx context.Context, id int) (*models.User, error) {
	panic(fmt.Errorf("not implemented: GetUser - GetUser"))
}

// GetCurrentUser is the resolver for the GetCurrentUser field.
func (r *queryResolver) GetCurrentUser(ctx context.Context) (*models.User, error) {
	panic(fmt.Errorf("not implemented: GetCurrentUser - GetCurrentUser"))
}

// GetUsers is the resolver for the GetUsers field.
func (r *queryResolver) GetUsers(ctx context.Context, pagination *models.Pagination, filter *models.UserFilter) (*models.GetUsersResponse, error) {
	panic(fmt.Errorf("not implemented: GetUsers - GetUsers"))
}

// GetFullUser is the resolver for the GetFullUser field.
func (r *queryResolver) GetFullUser(ctx context.Context, id int) (*models.FullUser, error) {
	panic(fmt.Errorf("not implemented: GetFullUser - GetFullUser"))
}

// GetRatingByAchs is the resolver for the GetRatingByAchs field.
func (r *queryResolver) GetRatingByAchs(ctx context.Context) (*models.RatingByAch, error) {
	panic(fmt.Errorf("not implemented: GetRatingByAchs - GetRatingByAchs"))
}

// GetRatingByStat is the resolver for the GetRatingByStat field.
func (r *queryResolver) GetRatingByStat(ctx context.Context, id int) (*models.RatingByStat, error) {
	panic(fmt.Errorf("not implemented: GetRatingByStat - GetRatingByStat"))
}

// GetAchievement is the resolver for the GetAchievement field.
func (r *queryResolver) GetAchievement(ctx context.Context, id int) (*models.Achievement, error) {
	panic(fmt.Errorf("not implemented: GetAchievement - GetAchievement"))
}

// GetAchievements is the resolver for the GetAchievements field.
func (r *queryResolver) GetAchievements(ctx context.Context, pagination *models.Pagination) (*models.GetAchievementsResponse, error) {
	panic(fmt.Errorf("not implemented: GetAchievements - GetAchievements"))
}

// GetEvent is the resolver for the GetEvent field.
func (r *queryResolver) GetEvent(ctx context.Context, id int) (*models.GetEvent, error) {
	panic(fmt.Errorf("not implemented: GetEvent - GetEvent"))
}

// GetEvents is the resolver for the GetEvents field.
func (r *queryResolver) GetEvents(ctx context.Context, pagination *models.Pagination) (*models.GetEventsResponse, error) {
	panic(fmt.Errorf("not implemented: GetEvents - GetEvents"))
}

// GetStat is the resolver for the GetStat field.
func (r *queryResolver) GetStat(ctx context.Context, id int) (*models.Stat, error) {
	panic(fmt.Errorf("not implemented: GetStat - GetStat"))
}

// GetStats is the resolver for the GetStats field.
func (r *queryResolver) GetStats(ctx context.Context, pagination *models.Pagination) (*models.GetStatsResponse, error) {
	panic(fmt.Errorf("not implemented: GetStats - GetStats"))
}

// Mutation returns server.MutationResolver implementation.
func (r *Resolver) Mutation() server.MutationResolver { return &mutationResolver{r} }

// Query returns server.QueryResolver implementation.
func (r *Resolver) Query() server.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
