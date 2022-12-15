package directives

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/pkg/utils"
)

func AdminOnly(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	claims, ok := utils.GetClaimsFromCtx(ctx)
	if !ok {
		return nil, &gqlerror.Error{
			Message: "access denied: no token in context",
		}
	}
	if claims.Role == models.SuperAdminRole || claims.Role == models.AdminRole {
		return next(ctx)
	}

	return nil, &gqlerror.Error{
		Message: "access denied: no required role",
	}

}
