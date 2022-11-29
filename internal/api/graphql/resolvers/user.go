package resolvers

import (
	"context"
	"database/sql"
	"gitlab.com/krespix/gamification-api/internal/models"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
	"gitlab.com/krespix/gamification-api/pkg/utils"
)

func (r *Resolver) GetUser(ctx context.Context, id int) (*apiModels.User, error) {
	user, err := r.userService.Get(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return &apiModels.User{
		ID:        int(user.ID),
		ForeignID: utils.SqlNullStringToString(user.ForeignID),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		DeletedAt: utils.SqlNullTimeToTime(user.DeletedAt),
		Role:      apiModels.Role(user.Role),
	}, nil
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
	return nil, nil
}
