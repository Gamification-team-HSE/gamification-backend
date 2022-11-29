package resolvers

import (
	"context"
)

func (r *Resolver) SendCode(ctx context.Context, email string) (interface{}, error) {
	err := r.authService.SendCode(ctx, email)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
	}, nil
}

func (r *Resolver) VerifyCode(ctx context.Context, email string, code int) (string, error) {
	return r.authService.VerifyCode(ctx, email, code)
}
