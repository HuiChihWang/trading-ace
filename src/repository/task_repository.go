package repository

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"time"
	"trading-ace/src/database"
	"trading-ace/src/model"
)

const tasksTableName = "tasks"

type SearchTasksCondition struct {
	UserID    string
	Type      model.TaskType
	Status    model.TaskStatus
	StartTime time.Time
	EndTime   time.Time
}

type TaskRepository interface {
	CreateTask(task *model.Task) (*model.Task, error)
	GetTaskByID(taskID int) (*model.Task, error)
	SearchTasks(condition *SearchTasksCondition) ([]*model.Task, error)
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

func (r *taskRepositoryImpl) SearchTasks(condition *SearchTasksCondition) ([]*model.Task, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := psql.Select("id, user_id, status, type, swap_amount, created_at, completed_at").From(tasksTableName)

	if condition.UserID != "" {
		query = query.Where(squirrel.Eq{"user_id": condition.UserID})
	}

	if condition.Type != "" {
		query = query.Where(squirrel.Eq{"type": condition.Type})
	}

	if condition.Status != "" {
		query = query.Where(squirrel.Eq{"status": condition.Status})
	}

	if !condition.StartTime.IsZero() || !condition.EndTime.IsZero() {

		if condition.StartTime.IsZero() || condition.EndTime.IsZero() {
			return nil, fmt.Errorf("both start time and end time should be provided")
		}

		if condition.StartTime.After(condition.EndTime) {
			return nil, fmt.Errorf("start time should be before end time")
		}

		query = query.Where(squirrel.Gt{"created_at": condition.StartTime.UTC()})
		query = query.Where(squirrel.Lt{"created_at": condition.EndTime.UTC()})
	}

	sqlCommand, args, err := query.OrderBy("id DESC").ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.dbInstance.Query(sqlCommand, args...)
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

func (r *taskRepositoryImpl) GetTaskByID(taskID int) (*model.Task, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlCommand, args, err := psql.
		Select("id, user_id, status, type, swap_amount, created_at, completed_at").
		From(tasksTableName).
		Where(squirrel.Eq{"id": taskID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := r.dbInstance.QueryRow(sqlCommand, args...)

	task := &model.Task{}
	err = row.Scan(&task.ID, &task.UserID, &task.Status, &task.Type, &task.SwapAmount, &task.CreatedAt, &task.CompletedAt)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *taskRepositoryImpl) CreateTask(task *model.Task) (*model.Task, error) {
	task.CreatedAt = task.CreatedAt.UTC()

	if task.CompletedAt.Valid {
		task.CompletedAt.Time = task.CompletedAt.Time.UTC()
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlCommand, args, err := psql.Insert(tasksTableName).
		Columns("user_id", "status", "type", "swap_amount", "created_at", "completed_at").
		Values(task.UserID, task.Status, task.Type, task.SwapAmount, task.CreatedAt, task.CompletedAt).
		Suffix("RETURNING id, user_id, status, type, swap_amount, created_at, completed_at").
		ToSql()

	err = r.dbInstance.QueryRow(sqlCommand, args...).Scan(&task.ID, &task.UserID, &task.Status, &task.Type, &task.SwapAmount, &task.CreatedAt, &task.CompletedAt)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *taskRepositoryImpl) UpdateTask(task *model.Task) (*model.Task, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query := psql.Update(tasksTableName)

	if task.Status != "" {
		query = query.Set("status", task.Status)
	}

	if task.CompletedAt.Valid {
		task.CompletedAt.Time = task.CompletedAt.Time.UTC()
		query = query.Set("completed_at", task.CompletedAt)
	}

	query = query.Where(squirrel.Eq{"id": task.ID}).
		Suffix("RETURNING status, completed_at")

	sqlCommand, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.dbInstance.QueryRow(sqlCommand, args...)
	err = row.Scan(&task.Status, &task.CompletedAt)

	if err != nil {
		return nil, err
	}

	return task, nil
}
