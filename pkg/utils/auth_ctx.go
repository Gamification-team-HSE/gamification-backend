package utils

import (
	"context"
	"gitlab.com/krespix/gamification-api/internal/models"
)

var (
	claimsContextKey struct{} = struct{}{}
)

func GetClaimsFromCtx(ctx context.Context) (*models.Claims, bool) {
	val, ok := ctx.Value(claimsContextKey).(*models.Claims)
	return val, ok
}

func PutClaimsToContext(ctx context.Context, claims *models.Claims) context.Context {
	ctx = context.WithValue(ctx, claimsContextKey, claims)
	return ctx
}
