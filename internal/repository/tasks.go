package repository

import (
	"flash_cards/internal"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type taskRepo struct {
	db *sqlx.DB
}

func (repo *taskRepo) Create() (internal.Task, error) {
	task := internal.Task{}
	var id int64
	if err := repo.db.QueryRow(
		"insert into flash_cards.tasks(details) values ('') returning id").Scan(&id); err != nil {
		return task, err
	}

	task, err := repo.Get(id)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (repo *taskRepo) Get(id int64) (internal.Task, error) {
	task := internal.Task{}
	sql, args, err := sq.Select(
		"id",
		"create_date",
		"status",
		"details").
		From("flash_cards.tasks").
		Where(sq.Eq{
			"id": id,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return task, err
	}
	tasks := []internal.Task{}
	if err := repo.db.Select(&tasks, sql, args...); err != nil {
		fmt.Println(sql, args)
		return task, err
	}
	if len(tasks) > 0 {
		task = tasks[0]
	}
	return task, nil
}

func (repo *taskRepo) SetStatus(id int64, status string) error {
	sql, args, err := sq.Update("flash_cards.tasks").
		Set("status", status).
		Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		fmt.Println("update builder error:", err)
		return err
	}
	if _, err := repo.db.Exec(sql, args...); err != nil {
		fmt.Println("update error:", err)
		return err
	}
	return nil
}
