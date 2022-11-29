package directives

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/krespix/gamification-api/pkg/utils"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	_, ok := utils.GetClaimsFromCtx(ctx)
	if !ok {
		return nil, &gqlerror.Error{
			Message: "access denied: no token in context",
		}
	}
	return next(ctx)
}
