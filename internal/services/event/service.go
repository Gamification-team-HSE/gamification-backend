package event

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"gitlab.com/krespix/gamification-api/internal/clients/s3"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/achievement"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/event"
	"gitlab.com/krespix/gamification-api/internal/services/image"
	apiModels "gitlab.com/krespix/gamification-api/pkg/graphql/models"
	"gitlab.com/krespix/gamification-api/pkg/utils"
	errors "gitlab.com/krespix/gamification-api/pkg/utils/graphq_erorrs"
)

type Service interface {
	Create(ctx context.Context, event *models.Event) error
	Update(ctx context.Context, event *models.UpdateEvent) error
	Get(ctx context.Context, id int) (*models.DbEvent, error)
	List(ctx context.Context, pagination *models.Pagination) (*apiModels.GetEventsResponse, error)
	Delete(ctx context.Context, id int) error
}
type service struct {
	validate  *validator.Validate
	eventRepo event.Repository
	achRepo   achievement.Repository

	folder   string
	s3Client s3.Client
}

func (s *service) Delete(ctx context.Context, id int) error {
	achList, err := s.achRepo.List(ctx, nil)
	if err != nil {
		return err
	}
	toUpdate := make([]*models.RepoAchievement, 0)

	needUpd := false
	for i, a := range achList {
		updRules := achList[i].Rules
		for j, b := range a.Rules.Blocks {
			res := make([]*models.EventRule, 0)
			for i, statRule := range b.EventRules {
				if statRule.EventID != id {
					res = append(res, b.EventRules[i])
				} else {
					needUpd = true
				}
			}
			updRules.Blocks[j].EventRules = res
		}
		if needUpd {
			toUpdate = append(toUpdate, &models.RepoAchievement{
				ID:    a.ID,
				Rules: updRules,
			})
		}
	}

	for _, a := range toUpdate {
		err = s.achRepo.Update(ctx, a)
		if err != nil {
			return err
		}
	}
	err = s.eventRepo.Delete(ctx, id)
	return err
}

func (s *service) List(ctx context.Context, pagination *models.Pagination) (*apiModels.GetEventsResponse, error) {
	var repoPagination *models.RepoPagination
	if pagination != nil {
		repoPagination = pagination.ToRepo()
	}
	events, err := s.eventRepo.List(ctx, repoPagination)
	if err != nil {
		return nil, err
	}
	total, err := s.eventRepo.Total(ctx)
	if err != nil {
		return nil, err
	}
	resEvents := make([]*apiModels.GetEvent, 0, len(events))
	for _, e := range events {
		ev := &apiModels.GetEvent{
			ID:          int(e.ID),
			Name:        e.Name,
			Description: utils.SqlNullStringToString(e.Description),
			CreatedAt:   int(e.CreatedAt.Unix()),
			StartAt:     int(e.StartAt.Unix()),
		}
		if e.Image.Valid {
			name := s.s3Client.BuildURL(s.folder, e.Image.String)
			ev.Image = &name
		}
		endAt := utils.SqlNullTimeToTime(e.EndAt)
		if endAt != 0 {
			ev.EndAt = &endAt
		}
		resEvents = append(resEvents, ev)
	}

	return &apiModels.GetEventsResponse{
		Total:  total,
		Events: resEvents,
	}, nil
}

func (s *service) Get(ctx context.Context, id int) (*models.DbEvent, error) {
	e, err := s.eventRepo.Get(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	if e.Image.Valid {
		e.Image.String = s.s3Client.BuildURL(s.folder, e.Image.String)
	}
	return e, nil
}

func (s *service) Create(ctx context.Context, event *models.Event) error {
	err := s.validate.Struct(event)
	if err != nil {
		return err
	}
	exists, err := s.eventRepo.ExistsByName(ctx, event.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("event with name %s already exists", event.Name)
	}

	if len(event.Description.String) > 1024 {
		return fmt.Errorf("event description exceed the maximum number of characters 1024")
	}

	if event.EndAt.Valid && event.EndAt.Time.Before(event.StartAt) {
		return fmt.Errorf("EndAt cannot be before StartAt")
	}

	if event.StartAt.Before(time.Now()) {
		fmt.Println("превышено")
		return fmt.Errorf("StartAt cannot be in the past")
	}

	if event.Image != nil {
		event.Image.Filename = image.GenerateFilename(event.Image)
		err = s.s3Client.Put(s.folder, event.Image.Filename, event.Image.ContentType, event.Image.File)
		if err != nil {
			return errors.InternalServerErrorWithDesc(ctx, err)
		}
	}

	return s.eventRepo.Create(ctx, event)

}

func (s *service) Update(ctx context.Context, event *models.UpdateEvent) error {
	err := s.validate.Struct(event)
	if err != nil {
		return errors.CustomError(ctx, 400, fmt.Sprintf("validation failed: %v", err))
	}
	updateEventReq := &models.UpdateEvent{ID: int64(event.ID)}

	if event.Image != nil {
		oldEvent, err := s.eventRepo.Get(ctx, int64(event.ID))
		if err != nil {
			fmt.Println("Тут")
			return errors.InternalServerErrorWithDesc(ctx, err)
		}
		if oldEvent.Image.String != "" {
			err = s.s3Client.Delete(s.folder, oldEvent.Image.String)
			if err != nil {
				return errors.InternalServerErrorWithDesc(ctx, err)
			}
		}

		event.Image.Filename = image.GenerateFilename(event.Image)
		updateEventReq.Image = event.Image

		err = s.s3Client.Put(s.folder, event.Image.Filename, event.Image.ContentType, event.Image.File)
		if err != nil {
			return errors.InternalServerErrorWithDesc(ctx, err)

		}
	}

	if event.Name.Valid {
		if len(event.Name.String) <= 128 {
			updateEventReq.Name = event.Name
		} else {
			return fmt.Errorf("event name exceed the maximum number of characters 128")
		}

	}
	if event.Description.Valid {
		if len(event.Description.String) <= 1024 {
			updateEventReq.Description = event.Description
		} else {
			return fmt.Errorf("event description exceed the maximum number of characters 1024")
		}

	}
	oldEvent, err := s.eventRepo.GetTime(ctx, int64(event.ID))
	if err != nil {
		return errors.InternalServerErrorWithDesc(ctx, err)
	}
	if event.EndAt.Valid {
		if event.StartAt.Valid && event.EndAt.Time.After(event.StartAt.Time) ||
			!event.StartAt.Valid && oldEvent.StartAt.Before(event.EndAt.Time) {
			updateEventReq.EndAt = event.EndAt
		} else {
			return fmt.Errorf("EndAt cannot be before StartAt")
		}

	}
	if event.StartAt.Valid && event.StartAt.Time.Before(oldEvent.CreatedAt) || !event.StartAt.Valid && oldEvent.StartAt.Before(oldEvent.CreatedAt) {
		return fmt.Errorf("StartAt cannot be in the past")
	} else {
		updateEventReq.StartAt = event.StartAt
	}

	return s.eventRepo.Update(ctx, event.ID, updateEventReq)
}

func New(
	eventRepo event.Repository,
	validate *validator.Validate,
	s3Client s3.Client,
	folder string,
	achRepo achievement.Repository,
) Service {
	return &service{
		eventRepo: eventRepo,
		validate:  validate,
		folder:    folder,
		s3Client:  s3Client,
		achRepo:   achRepo,
	}
}
