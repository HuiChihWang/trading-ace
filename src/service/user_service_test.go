package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"trading-ace/mock/repository"
	"trading-ace/src/model"
)

var userService UserService
var mockedUserRepository *repository.MockUserRepository

func setUpUserService(t *testing.T) {
	mockedUserRepository = repository.NewMockUserRepository(t)
	userService = &userServiceImpl{
		userRepository: mockedUserRepository,
	}
}

func TestUserServiceImpl_CreateUser(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		setUpUserService(t)

		mockedUserRepository.EXPECT().CreateUser("test_user_id").Return(
			&model.User{
				ID:     "test_user_id",
				Points: 0.0,
			}, nil).Times(1)

		newUser, err := userService.CreateUser("test_user_id")
		assert.Nil(t, err)
		assert.Equal(t, "test_user_id", newUser.ID)
		assert.Equal(t, 0.0, newUser.Points)
	})

	t.Run("CreateUserFail", func(t *testing.T) {
		setUpUserService(t)

		mockedUserRepository.EXPECT().CreateUser("test_user_id").Return(
			nil, assert.AnError).Times(1)

		newUser, err := userService.CreateUser("test_user_id")
		assert.Nil(t, newUser)
		assert.NotNil(t, err)
	})
}

func TestUserServiceImpl_GetUserByID(t *testing.T) {
	t.Run("GetUserByID", func(t *testing.T) {
		setUpUserService(t)

		mockedUserRepository.EXPECT().GetUser("test_user_id").Return(
			&model.User{
				ID:     "test_user_id",
				Points: 0.0,
			}, nil).Times(1)

		user, err := userService.GetUserByID("test_user_id")
		assert.Nil(t, err)
		assert.Equal(t, "test_user_id", user.ID)
		assert.Equal(t, 0.0, user.Points)
	})

	t.Run("GetUserByIDFail", func(t *testing.T) {
		setUpUserService(t)

		mockedUserRepository.EXPECT().GetUser("test_user_id").Return(
			nil, assert.AnError).Times(1)

		user, err := userService.GetUserByID("test_user_id")
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})
}

func TestUserServiceImpl_UpdateUserPoints(t *testing.T) {
	t.Run("UpdateUserPoints", func(t *testing.T) {
		setUpUserService(t)

		mockedUserRepository.EXPECT().GetUser("test_user_id").Return(
			&model.User{
				ID:     "test_user_id",
				Points: 0.0,
			}, nil).Times(1)

		mockedUserRepository.EXPECT().UpdateUser(&model.User{
			ID:     "test_user_id",
			Points: 100.0,
		}).Return(
			&model.User{
				ID:     "test_user_id",
				Points: 100.0,
			}, nil).Times(1)

		err := userService.UpdateUserPoints("test_user_id", 100.0)
		assert.Nil(t, err)
	})

	t.Run("UpdateUserWithNonPositivePoints", func(t *testing.T) {
		setUpUserService(t)

		t.Run("UpdateUserPointsWithZeroPoints", func(t *testing.T) {
			err := userService.UpdateUserPoints("test_user_id", 0.0)
			assert.NotNil(t, err)
		})

		t.Run("UpdateUserPointsWithNegativePoints", func(t *testing.T) {
			err := userService.UpdateUserPoints("test_user_id", -100.0)
			assert.NotNil(t, err)
		})
	})

	t.Run("CannotFindUser", func(t *testing.T) {
		setUpUserService(t)

		mockedUserRepository.EXPECT().GetUser("test_user_id").Return(
			nil, assert.AnError).Times(1)

		err := userService.UpdateUserPoints("test_user_id", 100.0)
		assert.NotNil(t, err)
	})
	t.Run("UpdateUserPointsFail", func(t *testing.T) {
		setUpUserService(t)

		mockedUserRepository.EXPECT().GetUser("test_user_id").Return(
			&model.User{
				ID:     "test_user_id",
				Points: 0.0,
			}, nil).Times(1)

		mockedUserRepository.EXPECT().UpdateUser(&model.User{
			ID:     "test_user_id",
			Points: 100.0,
		}).Return(
			nil, assert.AnError).Times(1)

		err := userService.UpdateUserPoints("test_user_id", 100.0)
		assert.NotNil(t, err)
	})
}
