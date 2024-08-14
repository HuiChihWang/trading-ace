package service

import (
	"errors"
	"trading-ace/src/model"
	"trading-ace/src/repository"
)

type UserService interface {
	GetUserByID(userID string) (*model.User, error)
	CreateUser(userID string) (*model.User, error)
	UpdateUserPoints(userID string, points float64) error
}

type userServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService() UserService {
	return &userServiceImpl{
		userRepository: repository.NewUserRepository(),
	}
}

func (s *userServiceImpl) CreateUser(userID string) (*model.User, error) {
	return s.userRepository.CreateUser(userID)
}

func (s *userServiceImpl) GetUserByID(userID string) (*model.User, error) {
	return s.userRepository.GetUser(userID)
}

func (s *userServiceImpl) UpdateUserPoints(userID string, points float64) error {
	if points <= 0 {
		return errors.New("points should be greater than 0")
	}

	user, err := s.userRepository.GetUser(userID)
	if err != nil {
		return err
	}

	user.Points += points

	_, err = s.userRepository.UpdateUser(user)
	return err
}
