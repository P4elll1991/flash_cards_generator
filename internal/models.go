package internal

import (
	"fmt"
	"reflect"
	"time"
)

type FlashCard struct {
	TaskId             int64  `json:"task_id" db:"task_id"`
	NativeLang         string `json:"native_lang" db:"native_lang"`
	LearningLang       string `json:"learning_lang" db:"learning_lang"`
	Topic              string `json:"topic" db:"topic"`
	Level              string `json:"level" db:"level"`
	Word               string `json:"word" db:"word"`
	Pronunciation      string `json:"pronunciation" db:"pronunciation"`
	PhoneticRespelling string `json:"phonetic_respelling" db:"phonetic_respelling"`
	Definition         string `json:"definition" db:"definition"`
	Translation        string `json:"translation" db:"translation"`
	Example            string `json:"example" db:"example"`
	ExampleTranslation string `json:"example_translation" db:"example_translation"`
}

func (card FlashCard) Hash() string {
	return fmt.Sprintf("%s-%s-%s-%s", card.NativeLang, card.LearningLang, card.Topic, card.Level)
}

type GenerateParams struct {
	TaskId       int64  `json:"-" db:"task_id"`
	NativeLang   string `json:"native_lang"  db:"native_lang"`
	LearningLang string `json:"learning_lang"  db:"learning_lang"`
	Words        int64  `json:"words"  db:"words"`
	Topic        string `json:"topic"  db:"topic"`
	Level        string `json:"level"  db:"level"`
}

func (param GenerateParams) Hash() string {
	return fmt.Sprintf("%s-%s-%s-%s", param.NativeLang, param.LearningLang, param.Topic, param.Level)
}

func (p GenerateParams) Divide(limit int64) []GenerateParams {
	params := []GenerateParams{}
	if p.Words <= limit {
		return []GenerateParams{p}
	}
	count := p.Words / limit
	if count > 0 {
		for i := 0; i < int(count); i++ {
			param := p
			param.Words = limit
			params = append(params, param)
		}
	}
	remainder := p.Words % limit
	if remainder > 0 {
		param := p
		param.Words = remainder
		params = append(params, param)
	}
	return params
}

func (params GenerateParams) GenerateRequest(mask string) string {
	return fmt.Sprintf(mask,
		params.NativeLang,
		params.LearningLang,
		params.Words,
		params.Topic,
		params.Level,
	)
}

func (params GenerateParams) fields() []interface{} {
	v := reflect.ValueOf(params)
	var fieldValues []interface{}
	for i := 0; i < v.NumField(); i++ {
		fieldValues = append(fieldValues, v.Field(i).Interface())
	}
	return fieldValues
}

type Task struct {
	Id          int64            `json:"id"  db:"id"`
	DateCreate  time.Time        `json:"create_date"  db:"create_date"`
	Progress    int64            `json:"progress"  db:"-"`
	Status      string           `json:"status"  db:"status"`
	Details     string           `json:"details"  db:"details"`
	RouteResult string           `json:"route_result"  db:"-"`
	Params      []GenerateParams `json:"params"  db:"-"`
}
