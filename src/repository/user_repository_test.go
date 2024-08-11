package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"trading-ace/src/config"
	"trading-ace/src/database"
)

func TestUserRepositoryImpl(t *testing.T) {
	setUpUserRepo := func(t *testing.T) *userRepositoryImpl {
		dbInstance := database.CreateDBInstance(&config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5435",
			Username: "postgres",
			Password: "postgres",
			DBName:   "trading_ace_test",
		})

		t.Cleanup(func() {
			dbInstance.Exec("DELETE FROM users")
			fmt.Println("[Tear Down] Cleaned up users table")
		})

		return &userRepositoryImpl{
			dbInstance: dbInstance,
		}
	}

	t.Run("CreateUser", func(t *testing.T) {
		repo := setUpUserRepo(t)
		user, err := repo.CreateUser("test_user_id")
		if err != nil {
			t.Errorf("CreateUser() error = %v", err)
		}

		assert.Equal(t, "test_user_id", user.ID)
		assert.Equal(t, 0.0, user.Points)
	})

	t.Run("GetUser", func(t *testing.T) {
		repo := setUpUserRepo(t)
		_, err := repo.CreateUser("test_user_id")
		if err != nil {
			t.Errorf("CreateUser() error = %v", err)
		}

		user, err := repo.GetUser("test_user_id")

		if err != nil {
			t.Errorf("GetUser() error = %v", err)
		}

		assert.Equal(t, "test_user_id", user.ID)
		assert.Equal(t, 0.0, user.Points)
	})
	t.Run("UpdateUser", func(t *testing.T) {
		repo := setUpUserRepo(t)
		user, err := repo.CreateUser("test_user_id")
		if err != nil {
			t.Errorf("CreateUser() error = %v", err)
		}

		user.Points = 100
		updatedUser, err := repo.UpdateUser(user)

		if err != nil {
			t.Errorf("UpdateUser() error = %v", err)
		}

		assert.Equal(t, "test_user_id", updatedUser.ID)
		assert.Equal(t, 100.0, updatedUser.Points)
	})
}
