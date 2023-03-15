package achievement

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres"
	"gitlab.com/krespix/gamification-api/pkg/utils"
)

const (
	achievementsTableName = "achievements"
)

type Repository interface {
	Create(ctx context.Context, achievement *models.RepoAchievement) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	Update(ctx context.Context, achievement *models.RepoAchievement) error
	Get(ctx context.Context, id int) (*models.RepoAchievement, error)
	Total(ctx context.Context) (int, error)
	List(ctx context.Context, achievement *models.RepoPagination) ([]*models.RepoAchievement, error)
	Delete(ctx context.Context, id int) error
}

type repository struct {
	*postgres.Client
}

func (r *repository) Delete(ctx context.Context, id int) error {
	qb := utils.PgQB().
		Delete(achievementsTableName).
		Where(sq.Eq{"id": id})
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.GetDBx().ExecContext(ctx, query, args...)
	return err
}

func (r *repository) List(ctx context.Context, pagination *models.RepoPagination) ([]*models.RepoAchievement, error) {
	qb := utils.PgQB().
		Select("*").
		From(achievementsTableName)
	if pagination != nil {
		qb = qb.Offset(uint64(pagination.Offset))
		qb = qb.Limit(uint64(pagination.Limit))
	}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	var res []*models.RepoAchievement
	err = r.GetDBx().SelectContext(ctx, &res, query, args...)
	return res, err
}

func (r *repository) Total(ctx context.Context) (int, error) {
	totalQuery := "select count(*) from achievements"
	var total int
	err := r.GetDBx().GetContext(ctx, &total, totalQuery)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *repository) Get(ctx context.Context, id int) (*models.RepoAchievement, error) {
	qb := utils.PgQB().
		Select("*").
		From(achievementsTableName).
		Where(sq.Eq{"id": id})
	res := &models.RepoAchievement{}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	err = r.GetDBx().GetContext(ctx, res, query, args...)
	fmt.Println(err)
	return res, err
}

func (r *repository) Update(ctx context.Context, achievement *models.RepoAchievement) error {
	qb := utils.PgQB().
		Update(achievementsTableName).
		Where(sq.Eq{"id": achievement.ID})
	if achievement.Name != "" {
		qb = qb.Set("name", achievement.Name)
	}
	if achievement.EndAt.Valid {
		qb = qb.Set("end_at", achievement.EndAt)
	}
	if achievement.Image.Valid {
		qb = qb.Set("image", achievement.Image)
	}
	if achievement.Rules != nil {
		qb = qb.Set("rules", achievement.Rules)
	}
	if achievement.Description.Valid {
		qb = qb.Set("description", achievement.Description)
	}
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.GetDBx().ExecContext(ctx, query, args...)
	return err
}

func (r *repository) ExistsByName(ctx context.Context, name string) (bool, error) {
	qb := utils.PgQB().Select("id").
		From(achievementsTableName).
		Where(sq.Eq{"name": name})
	query, args, err := qb.ToSql()
	if err != nil {
		return false, err
	}
	var id int64
	err = r.GetDBx().GetContext(ctx, &id, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *repository) Create(ctx context.Context, achievement *models.RepoAchievement) error {
	qb := utils.PgQB().
		Insert(achievementsTableName).
		Columns("name,"+
			"description,"+
			"image,"+
			"rules,"+
			"end_at,"+
			"created_at").
		Values(achievement.Name,
			achievement.Description,
			achievement.Image,
			achievement.Rules,
			achievement.EndAt,
			time.Now(),
		)

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = r.GetDBx().ExecContext(ctx, query, args...)

	return err
}

func New(client *postgres.Client) Repository {
	return &repository{client}
}
