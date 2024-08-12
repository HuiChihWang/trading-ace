package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"trading-ace/src/database"
	"trading-ace/src/exception"
	"trading-ace/src/model"
)

type UserRepository interface {
	CreateUser(id string) (*model.User, error)
	GetUser(id string) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
}

const usersTableName = "users"

type userRepositoryImpl struct {
	dbInstance *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{
		dbInstance: database.GetDBInstance(),
	}
}

func (u *userRepositoryImpl) CreateUser(id string) (*model.User, error) {
	user := model.NewUser(id)

	sqlCommand := fmt.Sprintf("INSERT INTO %s (id, points) VALUES ($1, $2)", usersTableName)
	stmt, err := u.dbInstance.Prepare(sqlCommand)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.Points)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepositoryImpl) GetUser(id string) (*model.User, error) {
	sqlCommand := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", usersTableName)
	row := u.dbInstance.QueryRow(sqlCommand, id)

	var user model.User
	err := row.Scan(&user.ID, &user.Points)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exception.UserNotFoundError
		}
		return nil, err
	}

	return &user, nil
}

func (u *userRepositoryImpl) UpdateUser(user *model.User) (*model.User, error) {
	sqlCommand := fmt.Sprintf("UPDATE %s SET points = $1 WHERE id = $2 RETURNING *", usersTableName)

	stmt, err := u.dbInstance.Prepare(sqlCommand)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.Points, user.ID).Scan(&user.ID, &user.Points)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exception.UserNotFoundError
		}
		return nil, err
	}

	return user, nil
}
