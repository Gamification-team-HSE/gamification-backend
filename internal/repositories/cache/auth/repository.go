package auth

import (
	"context"
	"fmt"
	"gitlab.com/krespix/gamification-api/internal/repositories/cache"
)

const (
	prefix = "auth_code"
)

type Repository interface {
	CreateCode(ctx context.Context, email string, code int) error
	GetCode(ctx context.Context, email string) (int, error)
}

type repository struct {
	client *cache.Client
}

func (r *repository) CreateCode(ctx context.Context, email string, code int) error {
	r.client.Set(key(email), code, 0)
	return nil
}

func (r *repository) GetCode(ctx context.Context, email string) (int, error) {
	res, found := r.client.Get(key(email))
	if !found {
		return 0, fmt.Errorf("not found")
	}
	code := res.(int)
	return code, nil
}

func key(key string) string {
	return fmt.Sprintf("%s_%s", prefix, key)
}

func New(client *cache.Client) Repository {
	return &repository{
		client: client,
	}
}
