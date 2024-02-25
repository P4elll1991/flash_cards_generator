package flashcards

import (
	"encoding/json"
	"errors"
	"flash_cards/internal"
	"fmt"
	"strings"
)

const openAI_request_mask = `Generate a list in json format using the following template:
[
	{
	 "topic"–  Asked Topic 
	 "level" -  Asked CEFR Level
	 "word"- the word on language, user wants to learn
	 "pronunciation" -  pronunciation instruction
	 "phonetic_respelling" - phonetic respelling using user's native language
	 "definition" - detailed word definition using user's native language
	 "translation" - translation of the word to user's native language,
	 "example" - example of usage of the word in real sentence.
	 "example_translation" -  translation of the usage example to user's native language
   }
]
For the following terms Native language: %s. Learning: %s. %d words, Topic: %s, CEFR Level: %s .
`

const exceptions_title = `Exclude the following words from the list: `

type OpenAIClient func(req string) (string, error)

var client OpenAIClient

func Init(c OpenAIClient) {
	client = c
}

func Generate(params internal.GenerateParams, exceptions map[string]internal.FlashCard) (map[string]internal.FlashCard, error) {
	if client == nil {
		return nil, errors.New("flash cards generator do not init")
	}

	req := params.GenerateRequest(openAI_request_mask)

	if len(exceptions) > 0 {
		req += exceptions_title
		words := []string{}
		for word := range exceptions {
			words = append(words, word)
		}
		req += strings.Join(words, ", ")
	}
	fmt.Println(req)
	resp, err := client(req)
	if err != nil {
		return nil, err
	}
	flashCards := []internal.FlashCard{}
	if err := json.Unmarshal([]byte(resp), &flashCards); err != nil {
		return nil, err
	}
	result := make(map[string]internal.FlashCard)
	for i := range flashCards {
		word := flashCards[i].Word
		if _, ok := exceptions[word]; !ok {
			flashCards[i].NativeLang = params.NativeLang
			flashCards[i].LearningLang = params.LearningLang
			result[word] = flashCards[i]
		}
	}

	return result, nil
}
