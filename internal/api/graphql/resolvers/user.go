package resolvers

import (
	"context"
	"database/sql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/krespix/gamification-api/internal/models"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
	"gitlab.com/krespix/gamification-api/pkg/utils"
)

func (r *Resolver) GetUser(ctx context.Context, id int) (*apiModels.User, error) {
	user, err := r.userService.Get(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return modelsUserToAPI(user), nil
}

func (r *Resolver) CreateUser(ctx context.Context, user apiModels.NewUser) (interface{}, error) {
	mUser := &models.User{
		Email: user.Email,
		Role:  models.Role(user.Role),
	}
	if user.ForeignID != nil {
		mUser.ForeignID = sql.NullString{
			String: *user.ForeignID,
			Valid:  true,
		}
	}
	if user.Name != nil {
		mUser.Name = sql.NullString{
			String: *user.Name,
			Valid:  true,
		}
	}
	err := r.userService.Create(ctx, mUser)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
	}, nil
}

func (r *Resolver) GetUsers(ctx context.Context) ([]*apiModels.User, error) {
	users, err := r.userService.List(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*apiModels.User, 0, len(users))
	for _, u := range users {
		res = append(res, modelsUserToAPI(u))
	}
	return res, nil
}

func (r *Resolver) GetCurrentUser(ctx context.Context) (*apiModels.User, error) {
	claims, ok := utils.GetClaimsFromCtx(ctx)
	if !ok {
		return nil, &gqlerror.Error{
			Message: "access denied: no token in context",
		}
	}
	user, err := r.userService.Get(ctx, claims.ID)
	if err != nil {
		return nil, err
	}
	return modelsUserToAPI(user), nil
}

func modelsUserToAPI(user *models.User) *apiModels.User {
	return &apiModels.User{
		ID:        int(user.ID),
		ForeignID: utils.SqlNullStringToString(user.ForeignID),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		DeletedAt: utils.SqlNullTimeToTime(user.DeletedAt),
		Role:      apiModels.Role(user.Role),
	}
}
