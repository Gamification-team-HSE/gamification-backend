package resolvers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/services/image"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
	"gitlab.com/krespix/gamification-api/pkg/utils"
	errors "gitlab.com/krespix/gamification-api/pkg/utils/graphq_erorrs"
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

func (r *Resolver) GetUsers(ctx context.Context, pagination *apiModels.Pagination, filter *apiModels.UserFilter) (*apiModels.GetUsersResponse, error) {
	var (
		mPagination *models.Pagination
		mFilter     *models.UserFilter
	)
	if pagination != nil {
		mPagination = &models.Pagination{
			Page: pagination.Page,
			Size: pagination.Size,
		}
	}
	if filter != nil {
		mFilter = &models.UserFilter{
			Active: boolPtrToVal(filter.Active),
			Banned: boolPtrToVal(filter.Banned),
			Admins: boolPtrToVal(filter.Admins),
		}
	}
	response, err := r.userService.List(ctx, mPagination, mFilter)
	if err != nil {
		return nil, err
	}
	userRes := make([]*apiModels.User, 0, len(response.Users))
	for _, u := range response.Users {
		userRes = append(userRes, modelsUserToAPI(u))
	}
	return &apiModels.GetUsersResponse{
		Users: userRes,
		Total: &apiModels.UsersTotalInfo{
			Admins: response.Total.Admins,
			Banned: response.Total.Banned,
			Active: response.Total.Active,
		},
	}, nil
}

func boolPtrToVal(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}

func (r *Resolver) UpdateUser(ctx context.Context, user *apiModels.UpdateUser) (interface{}, error) {
	claims, ok := utils.GetClaimsFromCtx(ctx)
	if !ok {
		return nil, &gqlerror.Error{
			Message: "access denied: no token in context",
		}
	}
	if claims.Role == models.DefaultRole {
		if int64(user.ID) != claims.ID {
			return nil, &gqlerror.Error{
				Message: "access denied: user can update only self profile",
			}
		}
	}

	if user.Avatar != nil {
		err := r.imageService.Validate(ctx, image.UserAvatarImage, user.Avatar)
		if err != nil {
			return nil, errors.CustomError(ctx, 400, fmt.Sprintf("image validation failed: %v", err))
		}
	}

	err := r.userService.Update(ctx, &models.UpdateUser{
		ID:     user.ID,
		Email:  stringPtrToValue(user.Email),
		Name:   stringPtrToValue(user.Name),
		Avatar: user.Avatar,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
	}, nil
}

func stringPtrToValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func (r *Resolver) DeleteUser(ctx context.Context, id int) (interface{}, error) {
	err := r.userService.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
	}, nil
}

func (r *Resolver) BanUser(ctx context.Context, id int) (interface{}, error) {
	err := r.userService.Ban(ctx, id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
	}, nil
}

func (r *Resolver) RecoverUser(ctx context.Context, id int) (interface{}, error) {
	err := r.userService.Recover(ctx, id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
	}, nil
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
	deletedAt := utils.SqlNullTimeToTime(user.DeletedAt)
	usr := &apiModels.User{
		ID:        int(user.ID),
		ForeignID: utils.SqlNullStringToString(user.ForeignID),
		Email:     user.Email,
		CreatedAt: int(user.CreatedAt.Unix()),
		Role:      apiModels.Role(user.Role),
		Name:      utils.SqlNullStringToString(user.Name),
		Avatar:    utils.SqlNullStringToString(user.Avatar),
	}
	if deletedAt != 0 {
		usr.DeletedAt = &deletedAt
	}
	return usr
}

func (r *Resolver) GetFullUser(ctx context.Context, id int) (*apiModels.FullUser, error) {
	fullUser, err := r.userService.GetFullUser(ctx, id)
	if err != nil {
		return nil, err
	}

	resStats := make([]*apiModels.UserStat, 0, len(fullUser.Stats))
	resEvents := make([]*apiModels.UserEvent, 0, len(fullUser.Events))
	resAch := make([]*apiModels.UserAch, 0, len(fullUser.Achievements))

	for _, ev := range fullUser.Events {
		resEvents = append(resEvents, &apiModels.UserEvent{
			EventID:     ev.EventID,
			Name:        ev.Name,
			Image:       utils.SqlNullStringToString(ev.Image),
			Description: utils.SqlNullStringToString(ev.Description),
			CreatedAt:   int(ev.CreatedAt.Unix()),
		})
	}
	for _, s := range fullUser.Stats {
		resStats = append(resStats, &apiModels.UserStat{
			StatID:      s.StatID,
			Name:        s.Name,
			Description: utils.SqlNullStringToString(s.Description),
			Value:       s.Value,
		})
	}

	for _, a := range fullUser.Achievements {
		resAch = append(resAch, &apiModels.UserAch{
			AchID:       a.AchID,
			Name:        a.Name,
			Description: utils.SqlNullStringToString(a.Description),
			CreatedAt:   int(a.CreatedAt.Unix()),
			Image:       utils.SqlNullStringToString(a.Image),
		})
	}

	res := &apiModels.FullUser{
		User:         modelsUserToAPI(fullUser.User),
		Stats:        resStats,
		Events:       resEvents,
		Achievements: resAch,
	}
	return res, nil
}
