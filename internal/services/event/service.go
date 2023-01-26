package event

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gitlab.com/krespix/gamification-api/internal/clients/s3"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/event"
	"gitlab.com/krespix/gamification-api/internal/services/image"
	errors "gitlab.com/krespix/gamification-api/pkg/utils/graphq_erorrs"
	"time"
)

type Service interface {
	Create(ctx context.Context, event *models.Event) error
	Update(ctx context.Context, event *models.UpdateEvent) error
}
type service struct {
	validate  *validator.Validate
	eventRepo event.Repository

	folder   string
	s3Client s3.Client
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

	//TODO run upload in goroutine
	if event.Image != nil {
		oldEvent, err := s.eventRepo.Get(ctx, int64(event.ID))
		if err != nil {
			fmt.Println("Тут")
			return errors.InternalServerErrorWithDesc(ctx, err)
		}
		if oldEvent.Image != "" {
			err = s.s3Client.Delete(s.folder, oldEvent.Image)
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
) Service {
	return &service{
		eventRepo: eventRepo,
		validate:  validate,
		folder:    folder,
		s3Client:  s3Client,
	}
}