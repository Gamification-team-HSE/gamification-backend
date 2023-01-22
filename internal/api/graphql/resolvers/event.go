package resolvers

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/services/image"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
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
	}

	if event.EndAt != nil {
		mEvent.EndAt = sql.NullTime{
			Time:  *event.EndAt,
			Valid: true,
		}
	}
	mEvent.StartAt = event.StartAt

	err := r.eventService.Create(ctx, mEvent)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status": "success",
	}, nil
}
