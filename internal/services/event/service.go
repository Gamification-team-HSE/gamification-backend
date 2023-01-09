package event

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/event"
)

type Service interface {
	Create(ctx context.Context, user *models.Event) error
}
type service struct {
	validate  *validator.Validate
	eventRepo event.Repository
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
	return s.eventRepo.Create(ctx, event)
}

func New(eventRepo event.Repository, validate *validator.Validate) Service {
	return &service{eventRepo: eventRepo, validate: validate}
}
