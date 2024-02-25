package repository

import (
	"flash_cards/internal"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type taskParamsRepo struct {
	db *sqlx.DB
}

func (repo *taskParamsRepo) Set(task_id int64, params []internal.GenerateParams) error {
	if len(params) > 0 {
		builder := sq.Insert("flash_cards.task_params").Columns(
			"task_id",
			"native_lang",
			"learning_lang",
			"topic",
			"level",
			"words",
		)

		for i := range params {
			builder = builder.Values(
				task_id,
				params[i].NativeLang,
				params[i].LearningLang,
				params[i].Topic,
				params[i].Level,
				params[i].Words,
			)
		}

		if _, err := builder.RunWith(repo.db).
			PlaceholderFormat(sq.Dollar).Exec(); err != nil {
			fmt.Println(builder.ToSql())
			return err
		}
	}
	return nil
}

func (repo *taskParamsRepo) Get(task_id int64) ([]internal.GenerateParams, error) {
	params := []internal.GenerateParams{}
	sql, args, err := sq.Select(
		"task_id",
		"native_lang",
		"learning_lang",
		"topic",
		"level",
		"words",
	).
		From("flash_cards.task_params").
		Where(sq.Eq{
			"task_id": task_id,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	if err := repo.db.Select(&params, sql, args...); err != nil {
		return nil, err
	}

	return params, nil
}

func (repo *taskParamsRepo) GetNewParams() ([]internal.GenerateParams, error) {
	tx, err := repo.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	result := []internal.GenerateParams{}
	sql, args, err := sq.Select(
		"p.task_id",
		"p.native_lang",
		"p.learning_lang",
		"p.topic",
		"p.level",
		"p.words",
	).
		From("flash_cards.task_params p").
		LeftJoin("flash_cards.tasks t on t.id = p.task_id").
		Where(sq.Eq{
			"t.status": "NEW",
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		fmt.Println("select builder error:", err)
		return nil, err
	}

	if err := tx.Select(&result, sql, args...); err != nil {
		fmt.Println("select error:", err)
		return nil, err
	}

	sql, args, err = sq.Update("flash_cards.tasks").
		Set("status", "PROCESSING").
		Where(sq.Eq{"status": "NEW"}).PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		fmt.Println("update builder error:", err)
		return nil, err
	}
	if _, err := tx.Exec(sql, args...); err != nil {
		fmt.Println("update error:", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("commit error:", err)
		return nil, err
	}
	return result, nil
}

func (repo *taskParamsRepo) GetUnfinishedTasks() ([]internal.GenerateParams, error) {
	result := []internal.GenerateParams{}
	sql, args, err := sq.Select(
		"p.task_id",
		"p.native_lang",
		"p.learning_lang",
		"p.topic",
		"p.level",
		"p.words",
	).
		From("flash_cards.task_params p").
		LeftJoin("flash_cards.tasks t on t.id = p.task_id").
		Where(sq.Eq{
			"t.status": "PROCESSING",
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		fmt.Println("select builder error:", err)
		return nil, err
	}

	if err := repo.db.Select(&result, sql, args...); err != nil {
		fmt.Println("select error:", err)
		return nil, err
	}

	return result, nil
}
