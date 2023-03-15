package achievement

import (
	"context"
	"database/sql"
	"fmt"

	"gitlab.com/krespix/gamification-api/internal/clients/s3"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres/achievement"
	"gitlab.com/krespix/gamification-api/internal/services/image"
)

type Service interface {
	Create(ctx context.Context, createAchievement *models.CreateAchievement) error
	Update(ctx context.Context, updateAchievement *models.UpdateAchievement) error
	Get(ctx context.Context, id int) (*models.RepoAchievement, error)
}

type service struct {
	achievementRepo achievement.Repository
	s3Client        s3.Client
	folder          string
}

func (s *service) Get(ctx context.Context, id int) (*models.RepoAchievement, error) {
	ach, err := s.achievementRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if ach.Image.Valid {
		ach.Image.String = s.s3Client.BuildURL(s.folder, ach.Image.String)
	}
	return ach, nil
}

func (s *service) Update(ctx context.Context, updateAchievement *models.UpdateAchievement) error {
	return nil
}

func (s *service) Create(ctx context.Context, createAchievement *models.CreateAchievement) error {
	exists, err := s.achievementRepo.ExistsByName(ctx, createAchievement.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("achievemnt with name \"%s\" already exists", createAchievement.Name)
	}
	ach := &models.RepoAchievement{
		Name:  createAchievement.Name,
		Rules: createAchievement.Rules,
	}
	if createAchievement.Description != "" {
		ach.Description = sql.NullString{
			String: createAchievement.Description,
			Valid:  true,
		}
	}
	if !createAchievement.EndAt.IsZero() {
		ach.EndAt = sql.NullTime{
			Time:  createAchievement.EndAt,
			Valid: true,
		}
	}
	if createAchievement.Image != nil {
		createAchievement.Image.Filename = image.GenerateFilename(createAchievement.Image)
		err := s.s3Client.Put(s.folder, createAchievement.Image.Filename, createAchievement.Image.ContentType, createAchievement.Image.File)
		if err != nil {
			return err
		}
		ach.Image = sql.NullString{
			String: createAchievement.Image.Filename,
			Valid:  true,
		}
	}
	err = s.achievementRepo.Create(ctx, ach)
	return err
}

func New(achievementRepo achievement.Repository, s3Client s3.Client, folder string) Service {
	return &service{achievementRepo: achievementRepo, s3Client: s3Client, folder: folder}
}
