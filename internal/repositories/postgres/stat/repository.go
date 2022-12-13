package stat

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres"
	"gitlab.com/krespix/gamification-api/pkg/utils"
	"time"
)

const (
	statsTableName = "stats"
)

type Repository interface {
	//Get получаем стату по id
	Get(ctx context.Context, id int64) (*models.Stat, error)
	// Create создание новой статы
	Create(ctx context.Context, stat *models.Stat) error
	// ExistsByName проверяет существует ли стата по наименованию
	ExistsByName(ctx context.Context, name string) (bool, error)
	// GetByName получение статы по наименованию
	GetByName(ctx context.Context, name string) (*models.Stat, error)
	List(ctx context.Context) ([]*models.Stat, error)
}

type repository struct {
	*postgres.Client
}

func (r *repository) List(ctx context.Context) ([]*models.Stat, error) {
	qb := utils.PgQB().Select("*").
		From(statsTableName)
	query, _, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	var stats []*models.Stat
	err = r.GetDBx().SelectContext(ctx, &stats, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return stats, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (*models.Stat, error) {
	qb := utils.PgQB().Select("*").
		From(statsTableName).
		Where(sq.Eq{"name": name})
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	stat := &models.Stat{}
	err = r.GetDBx().GetContext(ctx, stat, query, args...)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (r *repository) ExistsByName(ctx context.Context, name string) (bool, error) {
	qb := utils.PgQB().Select("id").
		From(statsTableName).
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

func (r *repository) Get(ctx context.Context, id int64) (*models.Stat, error) {
	qb := utils.PgQB().Select("*").
		From(statsTableName).
		Where(sq.Eq{"id": id})
	stat := &models.Stat{}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	err = r.GetDBx().GetContext(ctx, stat, query, args...)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (r *repository) Create(ctx context.Context, stat *models.Stat) error {
	qb := utils.PgQB().Insert(statsTableName).
		Columns(
			"name,"+
				"created_at,"+
				"start_at,"+
				"period").
		Values(
			stat.Name,
			time.Now(),
			stat.StartAt,
			stat.Period,
		)
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.GetDBx().ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func New(client *postgres.Client) Repository {
	return &repository{
		Client: client,
	}
}
