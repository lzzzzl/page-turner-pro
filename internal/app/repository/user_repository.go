package repository

import (
	"context"
	"strings"
	"time"

	"github.com/lzzzzl/page-turner-pro/internal/domain/common"
	"github.com/lzzzzl/page-turner-pro/internal/domain/model"
)

type repoUser struct {
	ID        int       `db:"id"`
	UID       string    `db:"string"`
	Email     string    `db:"string"`
	Name      string    `db:"string"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"update_at"`
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
	// insert := map[string]interface{}{
	// 	repoColumnUser.Name: param.Name,

	// }
	return nil, nil
}
