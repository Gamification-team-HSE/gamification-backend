package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gitlab.com/krespix/gamification-api/internal/core/config"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/user"
)

type Service interface {
	InitSuperAdmin(ctx context.Context, admin config.SuperAdmin) error
	Get(ctx context.Context, id int64) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

type service struct {
	validate *validator.Validate

	userRepo user.Repository
}

func (s *service) Create(ctx context.Context, user *models.User) error {
	err := s.validate.Struct(user)
	if err != nil {
		return err
	}
	exists, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}
	return s.userRepo.Create(ctx, user)
}

func (s *service) InitSuperAdmin(ctx context.Context, admin config.SuperAdmin) error {
	err := s.validate.Struct(admin)
	if err != nil {
		return err
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, admin.Email)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return s.userRepo.Create(ctx, &models.User{
		Email: admin.Email,
		Role:  models.SuperAdminRole,
		Name: sql.NullString{
			String: admin.Name,
			Valid:  true,
		},
	})
}

func (s *service) Get(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.Get(ctx, id)
}

func New(userRepo user.Repository, validate *validator.Validate) Service {
	return &service{userRepo: userRepo, validate: validate}
}
