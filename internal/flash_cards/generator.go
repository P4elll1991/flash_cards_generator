package flashcards

import (
	"encoding/json"
	"errors"
	"flash_cards/internal"
	"fmt"
	"os"
	"strings"
)

var (
	openAI_request              = os.Getenv("FLASH_CARDS_OPEN_AI_REQUEST")
	openAI_request_postfix_mask = os.Getenv("FLASH_CARDS_OPEN_AI_REQUEST_POSTFIX_MASK")
)

const exceptions_title = `Exclude the following words from the list: `

type OpenAIClient func(req string) (string, error)

var client OpenAIClient

func Init(c OpenAIClient) {
	fmt.Println(openAI_request + openAI_request_postfix_mask)
	client = c
}

func Generate(params internal.GenerateParams, exceptions map[string]internal.FlashCard) (map[string]internal.FlashCard, error) {
	if client == nil {
		return nil, errors.New("flash cards generator do not init")
	}

	req := params.GenerateRequest(openAI_request + openAI_request_postfix_mask)

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
