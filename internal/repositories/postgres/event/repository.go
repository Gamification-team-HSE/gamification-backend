package event

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
	eventsTableName     = "event"
	userEventsTableName = "user_events"
)

type Repository interface {
	// Create создание нового события
	Create(ctx context.Context, event *models.Event) error
	// ExistsByName проверяет существует ли событие по наименованию
	ExistsByName(ctx context.Context, name string) (bool, error)
	// Get получение события по id
	Get(ctx context.Context, id int64) (*models.DbEvent, error)
	// GetTime получение временных параметров по id
	GetTime(ctx context.Context, id int64) (*models.EventTime, error)
	// Update обновляет все поля, которые переданы в структуре
	Update(ctx context.Context, id int64, event *models.UpdateEvent) error
	List(ctx context.Context, pagination *models.RepoPagination) ([]*models.DbEvent, error)
	Total(ctx context.Context) (int, error)
	Delete(ctx context.Context, id int) error
	GetUserEvents(ctx context.Context, userID int) ([]*models.UserEvent, error)
}

type repository struct {
	*postgres.Client
}

func (r *repository) GetUserEvents(ctx context.Context, userID int) ([]*models.UserEvent, error) {
	qb := utils.PgQB().
		Select("e.id as event_id, " +
			"e.name as name," +
			"e.description as description," +
			"e.image as image," +
			"ue.created_at as created_at").
		From(fmt.Sprintf("%s as e", eventsTableName)).
		Join(fmt.Sprintf("%s as ue on ue.event_id = e.id", userEventsTableName)).
		Where(sq.Eq{"ue.user_id": userID})
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	var res []*models.UserEvent
	err = r.GetDBx().SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	qb := utils.PgQB().
		Delete(eventsTableName).
		Where(sq.Eq{"id": id})
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.GetDBx().ExecContext(ctx, query, args...)
	return err
}

func (r *repository) List(ctx context.Context, pagination *models.RepoPagination) ([]*models.DbEvent, error) {
	qb := utils.PgQB().Select("*").
		From(eventsTableName).
		OrderBy("created_at")
	if pagination != nil {
		qb = qb.Limit(uint64(pagination.Limit))
		qb = qb.Offset(uint64(pagination.Offset))
	}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	var res []*models.DbEvent
	err = r.GetDBx().SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repository) Total(ctx context.Context) (int, error) {
	totalQuery := "select count(*) from event"
	var total int
	err := r.GetDBx().GetContext(ctx, &total, totalQuery)
	if err != nil {
		return 0, err
	}
	return total, nil
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
	img := sql.NullString{}
	if event.Image != nil {
		img = sql.NullString{
			String: event.Image.Filename,
			Valid:  true,
		}
	}
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
			img,
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

func (r *repository) Update(ctx context.Context, id int64, event *models.UpdateEvent) error {
	qb := utils.PgQB().Update(eventsTableName).
		Where(sq.Eq{"id": id})

	if event != nil {
		if event.Image != nil {
			qb = qb.Set("image", event.Image.Filename)
		}
		if event.Name.Valid {
			qb = qb.Set("name", event.Name)
		}
		if event.Description.Valid {
			qb = qb.Set("description", event.Description.String)
		}
		if event.EndAt.Valid {
			qb = qb.Set("end_at", event.EndAt.Time)
		}

		if event.StartAt.Valid {
			qb = qb.Set("start_at", event.StartAt.Time)
		}

	}
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.GetDBx().Exec(query, args...)
	return err
}

func (r *repository) Get(ctx context.Context, id int64) (*models.DbEvent, error) {
	qb := utils.PgQB().Select("*").
		From(eventsTableName).
		Where(sq.Eq{"id": id})
	event := &models.DbEvent{}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.GetDBx().GetContext(ctx, event, query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return event, nil
}

func (r *repository) GetTime(ctx context.Context, id int64) (*models.EventTime, error) {
	qb := utils.PgQB().Select("created_at, start_at, end_at").
		From(eventsTableName).
		Where(sq.Eq{"id": id})
	event := &models.EventTime{}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.GetDBx().GetContext(ctx, event, query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return event, nil
}

func New(client *postgres.Client) Repository {
	return &repository{
		Client: client,
	}
}
