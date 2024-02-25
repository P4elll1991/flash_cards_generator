package repository

import (
	"os"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	DRIVER = os.Getenv("FLASH_CARDS_DB_DRIVER")
	DSN    = os.Getenv("FLASH_CARDS_DB_DSN")
)

type repository struct {
	Tasks      *taskRepo
	TaskParams *taskParamsRepo
	FlashCards *flashCardRepo
}

func New() (*repository, error) {
	db, err := sqlx.Open(DRIVER, DSN)
	if err != nil {
		return nil, err
	}
	return &repository{
		Tasks:      &taskRepo{db: db},
		TaskParams: &taskParamsRepo{db: db},
		FlashCards: &flashCardRepo{db: db},
	}, nil
}
