package resolvers

import (
	"gitlab.com/krespix/gamification-api/internal/services/auth"
	"gitlab.com/krespix/gamification-api/internal/services/event"
	"gitlab.com/krespix/gamification-api/internal/services/stat"
	"gitlab.com/krespix/gamification-api/internal/services/user"
	"gitlab.com/krespix/gamification-api/pkg/graphql/server"
)

// Mutation returns server.MutationResolver implementation.
func (r *Resolver) Mutation() server.MutationResolver { return r }

// Query returns server.QueryResolver implementation.
func (r *Resolver) Query() server.QueryResolver { return r }

type Resolver struct {
	userService  user.Service
	authService  auth.Service
	statService  stat.Service
	eventService event.Service
}

func New(userService user.Service, authService auth.Service, statService stat.Service, eventService event.Service) *Resolver {
	return &Resolver{
		userService:  userService,
		authService:  authService,
		statService:  statService,
		eventService: eventService,
	}
}
