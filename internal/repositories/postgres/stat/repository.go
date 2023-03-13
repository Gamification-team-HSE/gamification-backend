package stat

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres"
	"gitlab.com/krespix/gamification-api/pkg/utils"
)

const (
	statsTableName = "stats"
)

type Repository interface {
	// Create создание новой статы
	Create(ctx context.Context, stat *models.Stat) error
	// ExistsByName проверяет существует ли стата по наименованию
	ExistsByName(ctx context.Context, name string) (bool, error)
	Get(ctx context.Context, id int) (*models.Stat, error)
	List(ctx context.Context, pagination *models.RepoPagination) ([]*models.Stat, error)
	Delete(ctx context.Context, id int) error
	Total(ctx context.Context) (int, error)
	Update(ctx context.Context, stat *models.UpdateStat) error
}

type repository struct {
	*postgres.Client
}

func (r *repository) Update(ctx context.Context, stat *models.UpdateStat) error {
	qb := utils.PgQB().
		Update(statsTableName).
		Where(sq.Eq{"id": stat.ID})
	if stat.Name != "" {
		qb = qb.Set("name", stat.Name)
	}
	if stat.Description != "" {
		qb = qb.Set("description", stat.Description)
	}
	if !stat.StartedAt.IsZero() {
		qb = qb.Set("start_at", stat.StartedAt)
	}
	if stat.Period != "" {
		qb = qb.Set("period", stat.Period)
	}
	if stat.SeqPeriod != "" {
		qb = qb.Set("seq_period", stat.SeqPeriod)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = r.GetDBx().ExecContext(ctx, query, args...)
	return err
}

func (r *repository) List(ctx context.Context, pagination *models.RepoPagination) ([]*models.Stat, error) {
	qb := utils.PgQB().
		Select("*").
		From(statsTableName)
	if pagination != nil {
		qb = qb.Limit(uint64(pagination.Limit))
		qb = qb.Offset(uint64(pagination.Offset))
	}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	var stats []*models.Stat
	err = r.GetDBx().SelectContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *repository) Get(ctx context.Context, id int) (*models.Stat, error) {
	qb := utils.PgQB().
		Select("*").
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

func (r *repository) Total(ctx context.Context) (int, error) {
	totalQuery := "select count(*) from stats"
	var total int
	err := r.GetDBx().GetContext(ctx, &total, totalQuery)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	qb := utils.PgQB().
		Delete(statsTableName).
		Where(sq.Eq{"id": id})
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.GetDBx().ExecContext(ctx, query, args...)
	return err
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

func (r *repository) Create(ctx context.Context, stat *models.Stat) error {
	qb := utils.PgQB().Insert(statsTableName).
		Columns(
			"name,"+
				"description,"+
				"created_at,"+
				"start_at,"+
				"period,"+
				"seq_period").
		Values(
			stat.Name,
			stat.Description,
			time.Now(),
			stat.StartAt,
			stat.Period,
			stat.SeqPeriod,
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
