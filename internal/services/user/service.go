package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gitlab.com/krespix/gamification-api/internal/clients/s3"
	"gitlab.com/krespix/gamification-api/internal/core/config"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/user"
	"gitlab.com/krespix/gamification-api/internal/services/image"
	errors "gitlab.com/krespix/gamification-api/pkg/utils/graphq_erorrs"
)

type Service interface {
	InitSuperAdmin(ctx context.Context, admin config.SuperAdmin) error
	Get(ctx context.Context, id int64) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	List(ctx context.Context, pagination *models.Pagination, filter *models.UserFilter) (*models.GetUsersResponse, error)
	Ban(ctx context.Context, id int) error
	Recover(ctx context.Context, id int) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user *models.UpdateUser) error
}

type service struct {
	validate *validator.Validate

	userRepo user.Repository

	folder   string
	s3Client s3.Client
}

func (s *service) Update(ctx context.Context, user *models.UpdateUser) error {
	err := s.validate.Struct(user)
	if err != nil {
		return errors.CustomError(ctx, 400, fmt.Sprintf("validation failed: %v", err))
	}
	updateUserReq := &models.User{ID: int64(user.ID)}
	//TODO run upload in goroutine
	if user.Avatar != nil {
		oldUser, err := s.userRepo.Get(ctx, int64(user.ID))
		if err != nil {
			return errors.InternalServerErrorWithDesc(ctx, err)
		}

		if oldUser.Avatar.Valid {
			err = s.s3Client.Delete(s.folder, oldUser.Avatar.String)
			if err != nil {
				return errors.InternalServerErrorWithDesc(ctx, err)
			}
		}

		user.Avatar.Filename = image.GenerateFilename(user.Avatar)
		updateUserReq.Avatar = sql.NullString{
			String: user.Avatar.Filename,
			Valid:  true,
		}
		err = s.s3Client.Put(s.folder, user.Avatar.Filename, user.Avatar.ContentType, user.Avatar.File)
		if err != nil {
			return errors.InternalServerErrorWithDesc(ctx, err)
		}
	}
	if user.Name != "" {
		updateUserReq.Name = sql.NullString{
			String: user.Name,
			Valid:  true,
		}
	}
	if user.Email != "" {
		updateUserReq.Email = user.Email
	}
	return s.userRepo.Update(ctx, user.ID, updateUserReq)
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *service) Recover(ctx context.Context, id int) error {
	usr, err := s.userRepo.Get(ctx, int64(id))
	if err != nil {
		return errors.InternalServerError(ctx)
	}
	if !usr.DeletedAt.Valid {
		return errors.CustomError(ctx, 400, "bad request: user not banned")
	}
	return s.userRepo.Recover(ctx, id)
}

func (s *service) Ban(ctx context.Context, id int) error {
	usr, err := s.userRepo.Get(ctx, int64(id))
	if err != nil {
		return errors.InternalServerError(ctx)
	}
	if usr.DeletedAt.Valid {
		return errors.CustomError(ctx, 400, "bad request: user already banned")
	}
	return s.userRepo.SoftDelete(ctx, id)
}

func (s *service) List(ctx context.Context, pagination *models.Pagination, filter *models.UserFilter) (*models.GetUsersResponse, error) {
	var repoPagination *models.RepoPagination
	if pagination != nil {
		repoPagination = pagination.ToRepo()
	}
	users, err := s.userRepo.List(ctx, repoPagination, filter)
	if err != nil {
		return nil, err
	}
	totalInfo, err := s.userRepo.Total(ctx)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Avatar.Valid {
			u.Avatar.String = s.s3Client.BuildURL(s.folder, u.Avatar.String)
		}
	}

	return &models.GetUsersResponse{
		Users: users,
		Total: totalInfo,
	}, nil
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
	usr, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if usr.Avatar.Valid {
		usr.Avatar.String = s.s3Client.BuildURL(s.folder, usr.Avatar.String)
	}
	return usr, nil
}

func New(
	userRepo user.Repository,
	validate *validator.Validate,
	s3Client s3.Client,
	folder string,
) Service {
	return &service{
		userRepo: userRepo,
		validate: validate,
		folder:   folder,
		s3Client: s3Client,
	}
}
