package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

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

// CreateStat is the resolver for the CreateStat field.
func (r *mutationResolver) CreateStat(ctx context.Context, stat models.NewStat) (interface{}, error) {
	panic(fmt.Errorf("not implemented: CreateStat - CreateStat"))
}

// CreateEvent is the resolver for the CreateEvent field.
func (r *mutationResolver) CreateEvent(ctx context.Context, event models.NewEvent) (interface{}, error) {
	panic(fmt.Errorf("not implemented: CreateEvent - CreateEvent"))
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
func (r *queryResolver) GetUsers(ctx context.Context) ([]*models.User, error) {
	panic(fmt.Errorf("not implemented: GetUsers - GetUsers"))
}

// Mutation returns server.MutationResolver implementation.
func (r *Resolver) Mutation() server.MutationResolver { return &mutationResolver{r} }

// Query returns server.QueryResolver implementation.
func (r *Resolver) Query() server.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
