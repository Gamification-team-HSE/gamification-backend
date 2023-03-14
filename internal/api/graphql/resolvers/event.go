package resolvers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/services/image"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
	"gitlab.com/krespix/gamification-api/pkg/utils"
	errors "gitlab.com/krespix/gamification-api/pkg/utils/graphq_erorrs"
)

func (r *Resolver) CreateEvent(ctx context.Context, event apiModels.NewEvent) (interface{}, error) {
	mEvent := &models.Event{
		Name: event.Name,
	}
	if event.Description != nil {
		mEvent.Description = sql.NullString{
			String: *event.Description,
			Valid:  true,
		}
	}

	if event.Image != nil {
		err := r.imageService.Validate(ctx, image.EventIconImage, event.Image)
		if err != nil {
			return nil, errors.CustomError(ctx, 400, fmt.Sprintf("image validation failed: %v", err))
		}
		mEvent.Image = event.Image
	}

	if event.EndAt != nil {

		mEvent.EndAt = sql.NullTime{
			Time:  time.Unix(int64(*event.EndAt), 0),
			Valid: true,
		}
	}

	mEvent.StartAt = time.Unix(int64(event.StartAt), 0)

	err := r.eventService.Create(ctx, mEvent)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
	}, nil
}

func TimeToNullTime(ptime *time.Time) sql.NullTime {
	if ptime == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  *ptime,
		Valid: true,
	}
}

func stringToNullString(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *str,
		Valid:  true,
	}
}

func (r *Resolver) UpdateEvent(ctx context.Context, event apiModels.UpdateEvent) (interface{}, error) {
	_, ok := utils.GetClaimsFromCtx(ctx)
	if !ok {
		return nil, &gqlerror.Error{
			Message: "access denied: no token in context",
		}
	}

	if event.Image != nil {
		err := r.imageService.Validate(ctx, image.EventIconImage, event.Image)
		if err != nil {
			return nil, errors.CustomError(ctx, 400, fmt.Sprintf("image validation failed: %v", err))
		}
	}

	startAt := time.Unix(int64(*event.StartAt), 0)
	endAt := time.Unix(int64(*event.EndAt), 0)

	err := r.eventService.Update(ctx, &models.UpdateEvent{
		ID:          int64(event.ID),
		Name:        stringToNullString(event.Name),
		Description: stringToNullString(event.Description),
		Image:       event.Image,
		StartAt:     TimeToNullTime(&startAt),
		EndAt:       TimeToNullTime(&endAt),
	})
	if err != nil {
		fmt.Println("упал")
		return nil, err

	}
	return map[string]interface{}{
		"status": "success",
	}, nil
}

func (r *Resolver) GetEvent(ctx context.Context, id int) (*apiModels.GetEvent, error) {
	event, err := r.eventService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	endAt := utils.SqlNullTimeToTime(event.EndAt)
	res := &apiModels.GetEvent{
		ID:          int(event.ID),
		Name:        event.Name,
		Description: utils.SqlNullStringToString(event.Description),
		Image:       utils.SqlNullStringToString(event.Image),
		CreatedAt:   int(event.CreatedAt.Unix()),
		StartAt:     int(event.StartAt.Unix()),
	}
	if endAt != 0 {
		res.EndAt = &endAt
	}
	return nil, nil
}

func (r *Resolver) GetEvents(ctx context.Context, pagination *apiModels.Pagination) (*apiModels.GetEventsResponse, error) {
	var (
		mPagination *models.Pagination
	)
	if pagination != nil {
		mPagination = &models.Pagination{
			Page: pagination.Page,
			Size: pagination.Size,
		}
	}
	res, err := r.eventService.List(ctx, mPagination)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Resolver) DeleteEvent(ctx context.Context, id int) (interface{}, error) {
	err := r.eventService.Delete(ctx, id)
	return nil, err
}
