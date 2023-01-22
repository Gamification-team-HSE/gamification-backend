package event

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres"
	"gitlab.com/krespix/gamification-api/pkg/utils"
	"time"
)

const (
	eventsTableName = "event"
)

type Repository interface {
	// Create создание нового события
	Create(ctx context.Context, event *models.Event) error
	// ExistsByName проверяет существует ли событие по наименованию
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type repository struct {
	*postgres.Client
}

func (r *repository) ExistsByName(ctx context.Context, name string) (bool, error) {
	qb := utils.PgQB().Select("id").
		From(eventsTableName).
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

func (r *repository) Create(ctx context.Context, event *models.Event) error {
	qb := utils.PgQB().Insert(eventsTableName).
		Columns(
			"name,"+
				"description,"+
				"image,"+
				"created_at,"+
				"start_at,"+
				"end_at").
		Values(
			event.Name,
			event.Description,
			event.Image,
			time.Now(),
			event.StartAt,
			event.EndAt,
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

func (r *repository) UpdateImage(ctx context.Context, event *models.Event) error {
	qb := utils.PgQB().Update(eventsTableName).
		Set("image", event.Image).Where("id", event.ID)
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
