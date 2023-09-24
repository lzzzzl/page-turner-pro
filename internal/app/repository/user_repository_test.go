package repository

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/lzzzzl/page-turner-pro/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertUser(t *testing.T, expected *model.User, actual *model.User) {
	require.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.UID, actual.UID)
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, err := getPostgresDB()
	assert.NoError(t, err)

	repo := initRepository(t, db)

	// Args
	type Args struct {
		model.User
	}
	var args Args
	_ = faker.FakeData(&args)

	user, err := repo.CreateUser(context.Background(), args.User)
	require.NoError(t, err)
	assertUser(t, &args.User, user)
}
