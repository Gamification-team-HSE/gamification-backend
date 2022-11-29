package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

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

// GetUser is the resolver for the GetUser field.
func (r *queryResolver) GetUser(ctx context.Context, id int) (*models.User, error) {
	panic(fmt.Errorf("not implemented: GetUser - GetUser"))
}

// Mutation returns server.MutationResolver implementation.
func (r *Resolver) Mutation() server.MutationResolver { return &mutationResolver{r} }

// Query returns server.QueryResolver implementation.
func (r *Resolver) Query() server.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
