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
	// Create создание новой статы
	Create(ctx context.Context, stat *models.Stat) error
	// ExistsByName проверяет существует ли стата по наименованию
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type repository struct {
	*postgres.Client
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
