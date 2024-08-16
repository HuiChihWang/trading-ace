package repository

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"trading-ace/src/database"
	"trading-ace/src/exception"
	"trading-ace/src/model"
)

func TestUserRepositoryImpl(t *testing.T) {
	setUpUserRepo := func(t *testing.T) *userRepositoryImpl {
		dbInstance := database.GetDBInstance()
		t.Cleanup(func() {
			dbInstance.Exec("DELETE FROM users")
		})

		return &userRepositoryImpl{
			dbInstance: dbInstance,
		}
	}

	t.Run("CreateUser", func(t *testing.T) {
		repo := setUpUserRepo(t)
		user, err := repo.CreateUser("test_user_id")
		if err != nil {
			t.Errorf("CreateUser() exception = %v", err)
		}

		assert.Equal(t, "test_user_id", user.ID)
		assert.Equal(t, 0.0, user.Points)
	})

	t.Run("CreateUserDuplicate", func(t *testing.T) {
		repo := setUpUserRepo(t)
		_, err := repo.CreateUser("test_user_id")
		if err != nil {
			t.Errorf("CreateUser() exception = %v", err)
		}

		createdUser, err := repo.CreateUser("test_user_id")
		assert.Nil(t, createdUser)
		assert.NotNil(t, err)
	})

	t.Run("GetUser", func(t *testing.T) {
		repo := setUpUserRepo(t)
		_, err := repo.CreateUser("test_user_id")
		if err != nil {
			t.Errorf("CreateUser() exception = %v", err)
		}

		user, err := repo.GetUser("test_user_id")

		if err != nil {
			t.Errorf("GetUser() exception = %v", err)
		}

		assert.Equal(t, "test_user_id", user.ID)
		assert.Equal(t, 0.0, user.Points)
	})

	t.Run("GetUserNotFound", func(t *testing.T) {
		repo := setUpUserRepo(t)
		_, err := repo.CreateUser("test_user_id")
		if err != nil {
			t.Errorf("CreateUser() exception = %v", err)
		}

		_, err = repo.GetUser("not_found_user_id")
		assert.True(t, errors.Is(err, exception.UserNotFoundError))
	})

	t.Run("UpdateUser", func(t *testing.T) {
		repo := setUpUserRepo(t)
		user, err := repo.CreateUser("test_user_id")
		if err != nil {
			t.Errorf("CreateUser() exception = %v", err)
		}

		user.Points = 100
		updatedUser, err := repo.UpdateUser(user)

		if err != nil {
			t.Errorf("UpdateUser() exception = %v", err)
		}

		assert.Equal(t, "test_user_id", updatedUser.ID)
		assert.Equal(t, 100.0, updatedUser.Points)
	})

	t.Run("UpdateUserNotFound", func(t *testing.T) {
		repo := setUpUserRepo(t)
		_, err := repo.CreateUser("test_user_id")
		if err != nil {
			t.Errorf("CreateUser() exception = %v", err)
		}

		user := &model.User{
			ID:     "not_found_user_id",
			Points: 100,
		}

		_, err = repo.UpdateUser(user)
		assert.True(t, errors.Is(err, exception.UserNotFoundError))
	})
}
