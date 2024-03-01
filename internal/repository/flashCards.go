package repository

import (
	"flash_cards/internal"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type flashCardRepo struct {
	db *sqlx.DB
}

func (repo *flashCardRepo) CountWords(taskId int64) (int64, error) {
	result := []int64{}
	sql, args, err := sq.Select(
		"count(*)").
		From("flash_cards.cards").
		Where(sq.Eq{
			"task_id": taskId,
		}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	if err := repo.db.Select(&result, sql, args...); err != nil {
		return 0, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return 0, nil
}

func (repo *flashCardRepo) Search(params map[string]interface{}) ([]internal.FlashCard, error) {
	cards := []internal.FlashCard{}
	where := sq.Eq{}
	if task_id, ok := params["task_id"]; ok {
		where["task_id"] = task_id
	}

	if native_lang, ok := params["native_lang"]; ok {
		where["native_lang"] = native_lang
	}

	if learning_lang, ok := params["learning_lang"]; ok {
		where["learning_lang"] = learning_lang
	}

	if topic, ok := params["topic"]; ok {
		where["topic"] = topic
	}

	if level, ok := params["level"]; ok {
		where["level"] = level
	}

	if word, ok := params["word"]; ok {
		where["word"] = word
	}

	builder := sq.Select(
		"task_id",
		"word",
		"native_lang",
		"learning_lang",
		"topic",
		"level",
		"pronunciation",
		"phonetic_respelling",
		"definition",
		"translation",
		"example",
		"example_translation",
	).
		From("flash_cards.cards")

	if len(where) > 0 {
		builder = builder.Where(where)
	}

	sql, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	if err := repo.db.Select(&cards, sql, args...); err != nil {
		return nil, err
	}
	return cards, nil
}

func (repo *flashCardRepo) GetExceptions(params internal.GenerateParams) (map[string]internal.FlashCard, error) {
	searchParams := map[string]interface{}{
		"native_lang":   params.NativeLang,
		"learning_lang": params.LearningLang,
		"topic":         params.Topic,
	}

	cards, err := repo.Search(searchParams)
	if err != nil {
		return nil, err
	}

	result := make(map[string]internal.FlashCard, len(cards))
	for i := range cards {
		result[cards[i].Word] = cards[i]
	}

	return result, nil
}

func (repo *flashCardRepo) SaveCards(taskId int64, cards map[string]internal.FlashCard) error {
	if len(cards) > 0 {
		builder := sq.Insert("flash_cards.cards").Columns(
			"task_id",
			"word",
			"native_lang",
			"learning_lang",
			"topic",
			"level",
			"pronunciation",
			"phonetic_respelling",
			"definition",
			"translation",
			"example",
			"example_translation",
		)

		for i := range cards {
			builder = builder.Values(
				taskId,
				cards[i].Word,
				cards[i].NativeLang,
				cards[i].LearningLang,
				cards[i].Topic,
				cards[i].Level,
				cards[i].Pronunciation,
				cards[i].PhoneticRespelling,
				cards[i].Definition,
				cards[i].Translation,
				cards[i].Example,
				cards[i].ExampleTranslation,
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

func (repo *flashCardRepo) Update(cards []internal.FlashCard) error {
	if len(cards) > 0 {
		tx, err := repo.db.Beginx()
		if err != nil {
			return err
		}
		defer tx.Rollback()
		insertBuilder := sq.Insert("flash_cards.cards").Columns(
			"task_id",
			"word",
			"native_lang",
			"learning_lang",
			"topic",
			"level",
			"pronunciation",
			"phonetic_respelling",
			"definition",
			"translation",
			"example",
			"example_translation",
		)

		deletingCards := sq.Or{}

		for i := range cards {
			deletingCards = append(deletingCards, sq.Eq{
				"word":          cards[i].Word,
				"native_lang":   cards[i].NativeLang,
				"learning_lang": cards[i].LearningLang,
				"topic":         cards[i].Topic,
				"level":         cards[i].Level,
			})
			insertBuilder = insertBuilder.Values(
				cards[i].TaskId,
				cards[i].Word,
				cards[i].NativeLang,
				cards[i].LearningLang,
				cards[i].Topic,
				cards[i].Level,
				cards[i].Pronunciation,
				cards[i].PhoneticRespelling,
				cards[i].Definition,
				cards[i].Translation,
				cards[i].Example,
				cards[i].ExampleTranslation,
			)
		}
		if len(deletingCards) > 0 {
			if _, err := sq.Delete("flash_cards.cards").Where(deletingCards).
				RunWith(tx).PlaceholderFormat(sq.Dollar).Exec(); err != nil {
				return err
			}
		}

		if _, err := insertBuilder.RunWith(tx).
			PlaceholderFormat(sq.Dollar).Exec(); err != nil {
			fmt.Println(insertBuilder.ToSql())
			return err
		}

		if err = tx.Commit(); err != nil {
			fmt.Println("commit error:", err)
			return err
		}
	}
	return nil
}
