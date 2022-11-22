package resolvers

import (
	"context"
	"fmt"
	"gitlab.com/krespix/gamification-api/pkg/graphql/models"
	"time"
)

func (r *Resolver) GetUser(ctx context.Context, id int) (*models.User, error) {
	if id != 1 {
		return nil, fmt.Errorf("user not found")
	}
	return &models.User{
		ID:        1,
		ForeignID: 2,
		Email:     "krespix@gmail.com",
		CreatedAt: time.Now(),
		DeletedAt: nil,
		Role:      "admin",
	}, nil
}

func (r *Resolver) CreateUser(ctx context.Context, user models.NewUser) (*models.User, error) {
	return nil, nil
}
