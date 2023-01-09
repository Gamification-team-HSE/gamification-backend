package resolvers

import (
	"context"
	"database/sql"
	"gitlab.com/krespix/gamification-api/internal/models"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
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
		mEvent.Image = sql.NullString{
			String: *event.Image,
			Valid:  true,
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
