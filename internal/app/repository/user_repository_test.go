package repository

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/lzzzzl/page-turner-pro/internal/domain/model"
	"github.com/lzzzzl/page-turner-pro/testdata"
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
	db := getPostgresDB()
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

func TestUserRepository_GetUserByID(t *testing.T) {
	db := getPostgresDB()
	repo := initRepository(t, db, testdata.Path(testdata.TestDataUser))
	userID := 1

	_, err := repo.GetUserByID(context.Background(), userID)
	require.NoError(t, err)
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	db := getPostgresDB()
	repo := initRepository(t, db, testdata.Path(testdata.TestDataUser))
	email := "user1@pageturnerpro.com"

	_, err := repo.GetUserByEmail(context.Background(), email)
	require.NoError(t, err)
}

func TestUserRepository_GetAllUsers(t *testing.T) {
	db := getPostgresDB()
	repo := initRepository(t, db, testdata.Path(testdata.TestDataUser))

	users, err := repo.GetAllUsers(context.Background())
	require.NoError(t, err)
	assert.Len(t, users, 3)
}
