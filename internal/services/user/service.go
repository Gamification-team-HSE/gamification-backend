package user

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"github.com/go-playground/validator/v10"
	"gitlab.com/krespix/gamification-api/internal/clients/s3"
	"gitlab.com/krespix/gamification-api/internal/core/config"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/achievement"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/event"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/stat"
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
	GetFullUser(ctx context.Context, id int) (*models.FullUser, error)
	GetRatingByStat(ctx context.Context, statID int) (*models.RatingByStat, error)
	GetRatingByAchs(ctx context.Context) (*models.RatingByAchs, error)
	AddEvent(ctx context.Context, email string, eventID int) error
	AddStat(ctx context.Context, email string, statID int) error
}

type service struct {
	validate *validator.Validate

	userRepo         user.Repository
	statRepo         stat.Repository
	achievementsRepo achievement.Repository
	eventsRepo       event.Repository

	folder   string
	s3Client s3.Client
}

func (s *service) GetRatingByAchs(ctx context.Context) (*models.RatingByAchs, error) {
	users, err := s.userRepo.GetUserRatingByAchs(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].TotalAchs > users[j].TotalAchs
	})
	for i, u := range users {
		if u.Avatar.Valid {
			users[i].Avatar.String = s.s3Client.BuildURL(s.folder, u.Avatar.String)
		}
	}
	s.calculatePlacesByAchs(users)
	return &models.RatingByAchs{
		Total: len(users),
		Users: users,
	}, nil
}

func (s *service) calculatePlacesByAchs(users []*models.UserRatingByAch) {
	if len(users) == 0 {
		return
	}
	users[0].Place = 1
	prevVal := users[0].TotalAchs
	prevPlace := users[0].Place
	for i := 1; i < len(users); i++ {
		if prevVal == users[i].TotalAchs {
			users[i].Place = prevPlace
		} else {
			users[i].Place = i + 1
			prevPlace = users[i].Place
		}
		prevVal = users[i].TotalAchs
	}
}

func (s *service) GetRatingByStat(ctx context.Context, statID int) (*models.RatingByStat, error) {
	users, err := s.userRepo.GetUserRatingByStat(ctx, statID)
	if err != nil {
		return nil, err
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Value > users[j].Value
	})
	for i, u := range users {
		if u.Avatar.Valid {
			users[i].Avatar.String = s.s3Client.BuildURL(s.folder, u.Avatar.String)
		}
	}
	s.calculatePlacesByStat(users)
	return &models.RatingByStat{
		StatID: statID,
		Total:  len(users),
		Users:  users,
	}, nil
}

func (s *service) calculatePlacesByStat(users []*models.UserRatingByStat) {
	if len(users) == 0 {
		return
	}
	users[0].Place = 1
	prevVal := users[0].Value
	prevPlace := users[0].Place
	for i := 1; i < len(users); i++ {
		if prevVal == users[i].Value {
			users[i].Place = prevPlace
		} else {
			users[i].Place = i + 1
			prevPlace = users[i].Place
		}
		prevVal = users[i].Value
	}
}

func (s *service) GetFullUser(ctx context.Context, id int) (*models.FullUser, error) {
	userAchList, err := s.achievementsRepo.GetUsersAchievements(ctx, id)
	if err != nil {
		return nil, err
	}
	for i, v := range userAchList {
		if v.Image.Valid {
			userAchList[i].Image.String = s.s3Client.BuildURL("achievements", v.Image.String)
		}
	}
	userStats, err := s.statRepo.GetUserStats(ctx, id)
	if err != nil {
		return nil, err
	}
	userEvents, err := s.eventsRepo.GetUserEvents(ctx, id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for i, ue := range userEvents {
		if ue.Image.Valid {
			userEvents[i].Image.String = s.s3Client.BuildURL("events", ue.Image.String)
		}
	}
	u, err := s.userRepo.Get(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	if u.Avatar.Valid {
		u.Avatar.String = s.s3Client.BuildURL(s.folder, u.Avatar.String)
	}

	users, err := s.userRepo.GetUserRatingByAchs(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].TotalAchs > users[j].TotalAchs
	})
	s.calculatePlacesByAchs(users)
	place := -1
	for _, ru := range users {
		if ru.UserID == int(u.ID) {
			place = ru.Place
		}
	}

	return &models.FullUser{
		User:         u,
		Stats:        userStats,
		Events:       userEvents,
		Achievements: userAchList,
		Place:        place,
	}, nil
}

func (s *service) Update(ctx context.Context, user *models.UpdateUser) error {
	err := s.validate.Struct(user)
	if err != nil {
		return errors.CustomError(ctx, 400, fmt.Sprintf("validation failed: %v", err))
	}
	updateUserReq := &models.User{ID: int64(user.ID)}
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
	statRepo stat.Repository,
	achievementsRepo achievement.Repository,
	eventsRepo event.Repository,
) Service {
	return &service{
		userRepo:         userRepo,
		validate:         validate,
		folder:           folder,
		s3Client:         s3Client,
		statRepo:         statRepo,
		achievementsRepo: achievementsRepo,
		eventsRepo:       eventsRepo,
	}
}
