package resolvers

import (
	"gitlab.com/krespix/gamification-api/pkg/graphql/server"
)

// Mutation returns server.MutationResolver implementation.
func (r *Resolver) Mutation() server.MutationResolver { return r }

// Query returns server.QueryResolver implementation.
func (r *Resolver) Query() server.QueryResolver { return r }

type Resolver struct{}

func New() *Resolver {
	return &Resolver{}
}
