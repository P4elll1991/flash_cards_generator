export FLASH_CARDS_GENERATOR_CONCURRENCY=1000
export FLASH_CARDS_GENERATOR_LIMIT_WORDS_ONE_STEP=10
export FLASH_CARDS_GENERATOR_PORT=8000
export FLASH_CARDS_DB_DRIVER=pgx
export FLASH_CARDS_DB_DSN="host=localhost port=5432 user=postgres password=123 database=test"
go run .