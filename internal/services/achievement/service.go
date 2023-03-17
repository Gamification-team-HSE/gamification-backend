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
	List(ctx context.Context, pagination *models.Pagination) (*models.GetAchievementsResponse, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	achievementRepo achievement.Repository
	s3Client        s3.Client
	folder          string
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.achievementRepo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, pagination *models.Pagination) (*models.GetAchievementsResponse, error) {
	total, err := s.achievementRepo.Total(ctx)
	if err != nil {
		return nil, err
	}
	var repoPagination *models.RepoPagination
	if pagination != nil {
		repoPagination = pagination.ToRepo()
	}
	list, err := s.achievementRepo.List(ctx, repoPagination)
	if err != nil {
		return nil, err
	}

	for i, item := range list {
		if item.Image.Valid {
			list[i].Image.String = s.s3Client.BuildURL(s.folder, item.Image.String)
		}
	}

	return &models.GetAchievementsResponse{
		Achievements: list,
		Total:        total,
	}, nil
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
	old, err := s.Get(ctx, updateAchievement.ID)
	if err != nil {
		return err
	}
	repoAch := &models.RepoAchievement{
		ID: updateAchievement.ID,
	}
	if updateAchievement.Image != nil {
		if old.Image.Valid {
			err = s.s3Client.Delete(s.folder, old.Image.String)
			if err != nil {
				return err
			}
		}
		updateAchievement.Image.Filename = image.GenerateFilename(updateAchievement.Image)
		err = s.s3Client.Put(s.folder, updateAchievement.Image.Filename, updateAchievement.Image.ContentType, updateAchievement.Image.File)
		if err != nil {
			return err
		}
		repoAch.Image = sql.NullString{
			String: updateAchievement.Image.Filename,
			Valid:  true,
		}
	}
	if updateAchievement.Name != "" {
		repoAch.Name = updateAchievement.Name
	}
	if updateAchievement.Description != "" {
		repoAch.Description = sql.NullString{
			String: updateAchievement.Description,
			Valid:  true,
		}
	}
	if updateAchievement.Rules != nil {
		repoAch.Rules = updateAchievement.Rules
	}
	if !updateAchievement.EndAt.IsZero() {
		repoAch.EndAt = sql.NullTime{
			Time:  updateAchievement.EndAt,
			Valid: true,
		}
	}

	return s.achievementRepo.Update(ctx, repoAch)
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
