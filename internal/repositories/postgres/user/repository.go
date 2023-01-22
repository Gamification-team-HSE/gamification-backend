package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres"
	"gitlab.com/krespix/gamification-api/pkg/utils"
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
	// List пагинированный список юзеров с фильтром
	List(ctx context.Context, pagination *models.RepoPagination, filter *models.UserFilter) ([]*models.User, error)
	// Total подсчет категорий пользователей (активные, забаненные и админы)
	Total(ctx context.Context) (*models.UsersTotalInfo, error)
	// Delete навсегда удаляет запись
	Delete(ctx context.Context, id int) error
	// SoftDelete проставляет deleted_at в time.Now()
	SoftDelete(ctx context.Context, id int) error
	// Update обновляет все поля, которые переданы в структуре
	Update(ctx context.Context, id int, user *models.User) error
	// Recover очищает поле deleted_at
	Recover(ctx context.Context, id int) error
}

type repository struct {
	*postgres.Client
}

func (r *repository) Recover(ctx context.Context, id int) error {
	qb := utils.PgQB().Update(usersTableName).
		Where(sq.Eq{"id": id}).
		Set("deleted_at", nil)
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.GetDBx().Exec(query, args...)
	return err
}

func (r *repository) Update(ctx context.Context, id int, user *models.User) error {
	qb := utils.PgQB().Update(usersTableName).
		Where(sq.Eq{"id": id})
	if user != nil {
		if user.Avatar.Valid {
			qb = qb.Set("avatar", user.Avatar)
		}
		if user.Name.Valid {
			qb = qb.Set("name", user.Name)
		}
		if user.Email != "" {
			qb = qb.Set("email", user.Email)
		}
	}
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.GetDBx().Exec(query, args...)
	return err
}

func (r *repository) SoftDelete(ctx context.Context, id int) error {
	qb := utils.PgQB().Update(usersTableName).
		Where(sq.Eq{"id": id}).
		Set("deleted_at", time.Now())
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.GetDBx().Exec(query, args...)
	return err
}

func (r *repository) Delete(ctx context.Context, id int) error {
	qb := utils.PgQB().Delete(usersTableName).
		Where(sq.Eq{"id": id})
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.GetDBx().Exec(query, args...)
	return err
}

func (r *repository) Total(ctx context.Context) (*models.UsersTotalInfo, error) {
	adminsTotalQuery := "select count(id) from users where (role = 'admin' or role = 'super_admin') and deleted_at is null"
	activeTotalQuery := "select count(id) from users where deleted_at is null and role=user"
	bannedTotalQuery := "select count(id) from users where deleted_at is not null"
	query := fmt.Sprintf("select (%s) as admins, (%s) as active, (%s) as banned", adminsTotalQuery, activeTotalQuery, bannedTotalQuery)
	totalInfo := &models.UsersTotalInfo{}
	err := r.GetDBx().GetContext(ctx, totalInfo, query)
	if err != nil {
		return nil, err
	}
	return totalInfo, nil
}

func (r *repository) List(ctx context.Context, pagination *models.RepoPagination, filter *models.UserFilter) ([]*models.User, error) {
	qb := utils.PgQB().Select("*").
		From(usersTableName).
		OrderBy("created_at")
	if pagination != nil {
		qb = qb.Limit(uint64(pagination.Limit))
		qb = qb.Offset(uint64(pagination.Offset))
	}
	if filter != nil {
		if filter.Active {
			qb = qb.Where(sq.Expr("deleted_at is null and role=user"))
		}
		if filter.Banned {
			qb = qb.Where(sq.Expr("deleted_at is not null"))
		}
		if filter.Admins {
			qb = qb.Where(sq.And{
				sq.Or{
					sq.Eq{"role": models.SuperAdminRole},
					sq.Eq{"role": models.AdminRole},
				},
				sq.Expr("deleted_at is null"),
			})
		}
	}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	var users []*models.User
	err = r.GetDBx().SelectContext(ctx, &users, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return users, nil
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
