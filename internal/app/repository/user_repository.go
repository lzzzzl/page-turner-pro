package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lzzzzl/page-turner-pro/internal/domain/common"
	"github.com/lzzzzl/page-turner-pro/internal/domain/model"
)

type repoUser struct {
	ID        int       `db:"id"`
	UID       string    `db:"uid"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type repoColumnPatternUser struct {
	ID        string
	UID       string
	Email     string
	Name      string
	CreatedAt string
	UpdatedAt string
}

const repoTableUser = "users"

var repoColumnUser = repoColumnPatternUser{
	ID:        "id",
	UID:       "uid",
	Email:     "email",
	Name:      "name",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

func (c *repoColumnPatternUser) columns() string {
	return strings.Join([]string{
		c.ID,
		c.UID,
		c.Email,
		c.Name,
		c.CreatedAt,
		c.UpdatedAt,
	}, ", ")
}

func (r *PostgresRepository) CreateUser(ctx context.Context, param model.User) (*model.User, common.Error) {
	insert := map[string]interface{}{
		repoColumnUser.Name:  param.Name,
		repoColumnUser.UID:   param.UID,
		repoColumnUser.Email: param.Email,
		repoColumnUser.Name:  param.Name,
	}

	// build SQL query
	query, args, err := r.pgsq.Insert(repoTableUser).
		SetMap(insert).
		Suffix(fmt.Sprintf("returning %s", repoColumnUser.columns())).
		ToSql()
	if err != nil {
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}

	// execute SQL query
	var row repoUser
	if err = r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}

	user := model.User(row)

	return &user, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id int) (*model.User, common.Error) {
	where := sq.And{
		sq.Eq{repoColumnUser.ID: id},
	}

	// build SQL query
	query, args, err := r.pgsq.Select(repoColumnUser.columns()).
		From(repoTableUser).
		Where(where).
		Limit(1).
		ToSql()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.NewError(common.ErrorCodeResourceNotFound, err)
		}
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}

	// execute SQL query
	var row repoUser
	if err = r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}

	user := model.User(row)
	return &user, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, common.Error) {
	where := sq.And{
		sq.Eq{repoColumnUser.Email: email},
	}

	// build SQL query
	query, args, err := r.pgsq.Select(repoColumnUser.columns()).
		From(repoTableUser).
		Where(where).
		Limit(1).
		ToSql()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.NewError(common.ErrorCodeResourceNotFound, err)
		}
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}

	// execute SQL query
	var row repoUser
	if err = r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}

	user := model.User(row)
	return &user, nil
}

func (r *PostgresRepository) GetAllUsers(ctx context.Context) ([]*model.User, common.Error) {
	where := sq.And{}

	// build SQL query
	query, args, err := r.pgsq.Select(repoColumnUser.columns()).
		From(repoTableUser).
		Where(where).
		OrderBy(fmt.Sprintf("%s desc", repoColumnUser.CreatedAt)).
		ToSql()
	if err != nil {
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}

	// execute SQL query
	var rows []repoUser
	if err = r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}

	var users []*model.User
	for _, row := range rows {
		user := model.User(row)
		users = append(users, &user)
	}

	return users, nil
}

// func (r *PostgresRepository) UpdateUser(ctx context.Context, param model.User) ()
