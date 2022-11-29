package user

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
	usersTableName = "users"
)

type Repository interface {
	// Get получаем юзера по id
	Get(ctx context.Context, id int64) (*models.User, error)
	// Create создание нового юзера
	Create(ctx context.Context, user *models.User) error
	// ExistsByEmail проверяет существует ли юзер по емейл
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	// GetByEmail получение юезра по емейлу
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type repository struct {
	*postgres.Client
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	qb := utils.PgQB().Select("*").
		From(usersTableName).
		Where(sq.Eq{"email": email})
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	user := &models.User{}
	err = r.GetDBx().GetContext(ctx, user, query, args...)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	qb := utils.PgQB().Select("id").
		From(usersTableName).
		Where(sq.Eq{"email": email})
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

func (r *repository) Get(ctx context.Context, id int64) (*models.User, error) {
	qb := utils.PgQB().Select("*").
		From(usersTableName).
		Where(sq.Eq{"id": id})
	user := &models.User{}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	err = r.GetDBx().GetContext(ctx, user, query, args...)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) Create(ctx context.Context, user *models.User) error {
	qb := utils.PgQB().Insert(usersTableName).
		Columns(
			"foreign_id,"+
				"email,"+
				"created_at,"+
				"role,"+
				"name").
		Values(
			user.ForeignID,
			user.Email,
			time.Now(),
			user.Role,
			user.Name,
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
