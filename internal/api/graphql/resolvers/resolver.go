package resolvers

import (
	"gitlab.com/krespix/gamification-api/internal/services/achievement"
	"gitlab.com/krespix/gamification-api/internal/services/auth"
	"gitlab.com/krespix/gamification-api/internal/services/event"
	"gitlab.com/krespix/gamification-api/internal/services/image"
	"gitlab.com/krespix/gamification-api/internal/services/stat"
	"gitlab.com/krespix/gamification-api/internal/services/user"
	"gitlab.com/krespix/gamification-api/pkg/graphql/server"
)

// Mutation returns server.MutationResolver implementation.
func (r *Resolver) Mutation() server.MutationResolver { return r }

// Query returns server.QueryResolver implementation.
func (r *Resolver) Query() server.QueryResolver { return r }

type Resolver struct {
	userService        user.Service
	statService        stat.Service
	authService        auth.Service
	imageService       image.Service
	eventService       event.Service
	achievementService achievement.Service
}

func New(
	userService user.Service,
	authService auth.Service,
	statService stat.Service,
	imgService image.Service,
	eventService event.Service,
	achievementService achievement.Service,
) *Resolver {
	return &Resolver{
		userService:        userService,
		authService:        authService,
		statService:        statService,
		imageService:       imgService,
		eventService:       eventService,
		achievementService: achievementService,
	}
}
