package repository

import (
	"database/sql"
	"fmt"
	"time"
	"trading-ace/src/database"
	"trading-ace/src/model"
)

const tasksTableName = "tasks"

type TaskRepository interface {
	CreateTask(task *model.Task) (*model.Task, error)
	GetTasksByDateRange(from time.Time, to time.Time) ([]*model.Task, error)
	GetTaskByID(taskID int) (*model.Task, error)
	GetTasksByUserID(userID string) ([]*model.Task, error)
	GetTasksByUserIDAndType(userID string, taskType model.TaskType) ([]*model.Task, error)
	UpdateTask(task *model.Task) (*model.Task, error)
}

type taskRepositoryImpl struct {
	dbInstance *sql.DB
}

func NewTaskRepository() TaskRepository {
	return &taskRepositoryImpl{
		dbInstance: database.GetDBInstance(),
	}
}

func (r *taskRepositoryImpl) CreateTask(task *model.Task) (*model.Task, error) {
	if task.CompletedAt.Valid {
		task.CompletedAt.Time = task.CompletedAt.Time.In(time.UTC)
	}
	task.CreatedAt = task.CreatedAt.In(time.UTC)

	sqlCommand := fmt.Sprintf("INSERT INTO %s (user_id, status, type, swap_amount, created_at, completed_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", tasksTableName)
	stmt, err := r.dbInstance.Prepare(sqlCommand)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(task.UserID, task.Status, task.Type, task.SwapAmount, task.CreatedAt, task.CompletedAt).Scan(&task.ID)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *taskRepositoryImpl) UpdateTask(task *model.Task) (*model.Task, error) {
	if task.CompletedAt.Valid {
		task.CompletedAt.Time = task.CompletedAt.Time.In(time.UTC)
	}

	sqlCommand := fmt.Sprintf("UPDATE %s SET status = $1, completed_at = $2 WHERE id = $3", tasksTableName)

	stmt, err := r.dbInstance.Prepare(sqlCommand)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Status, task.CompletedAt, task.ID)

	if err != nil {
		return nil, err
	}

	updatedTask, err := r.GetTaskByID(task.ID)

	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

func (r *taskRepositoryImpl) GetTaskByID(taskID int) (*model.Task, error) {
	sqlCommand := fmt.Sprintf("SELECT id, user_id, status, type, swap_amount, created_at, completed_at FROM %s WHERE id = $1", tasksTableName)
	row := r.dbInstance.QueryRow(sqlCommand, taskID)

	task := &model.Task{}
	err := row.Scan(&task.ID, &task.UserID, &task.Status, &task.Type, &task.SwapAmount, &task.CreatedAt, &task.CompletedAt)
	if err != nil {
		return nil, err
	}

	task.CreatedAt = task.CreatedAt.In(time.UTC)
	if task.CompletedAt.Valid {
		task.CompletedAt.Time = task.CompletedAt.Time.In(time.UTC)
	}

	return task, nil
}

func (r *taskRepositoryImpl) GetTasksByDateRange(from time.Time, to time.Time) ([]*model.Task, error) {
	sqlCommand := fmt.Sprintf("SELECT id, user_id, status, type, swap_amount, created_at, completed_at FROM %s WHERE created_at >= $1 AND created_at < $2", tasksTableName)
	rows, err := r.dbInstance.Query(sqlCommand, from, to)

	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Status, &task.Type, &task.SwapAmount, &task.CreatedAt, &task.CompletedAt)
		if err != nil {
			return nil, err
		}

		task.CreatedAt = task.CreatedAt.In(time.UTC)
		if task.CompletedAt.Valid {
			task.CompletedAt.Time = task.CompletedAt.Time.In(time.UTC)
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *taskRepositoryImpl) GetTasksByUserIDAndType(userID string, taskType model.TaskType) ([]*model.Task, error) {
	sqlCommand := fmt.Sprintf("SELECT id, user_id, status, type, swap_amount, created_at, completed_at FROM %s WHERE user_id = $1 AND type = $2", tasksTableName)
	rows, err := r.dbInstance.Query(sqlCommand, userID, taskType)

	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Status, &task.Type, &task.SwapAmount, &task.CreatedAt, &task.CompletedAt)
		if err != nil {
			return nil, err
		}

		task.CreatedAt = task.CreatedAt.In(time.UTC)
		if task.CompletedAt.Valid {
			task.CompletedAt.Time = task.CompletedAt.Time.In(time.UTC)
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *taskRepositoryImpl) GetTasksByUserID(userID string) ([]*model.Task, error) {
	sqlCommand := fmt.Sprintf("SELECT id, user_id, status, type, swap_amount, created_at, completed_at FROM %s WHERE user_id = $1", tasksTableName)
	rows, err := r.dbInstance.Query(sqlCommand, userID)

	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Status, &task.Type, &task.SwapAmount, &task.CreatedAt, &task.CompletedAt)

		if err != nil {
			return nil, err
		}

		task.CreatedAt = task.CreatedAt.In(time.UTC)
		if task.CompletedAt.Valid {
			task.CompletedAt.Time = task.CompletedAt.Time.In(time.UTC)
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}
