package openai

import (
	"context"
	"log"
	"os"

	"github.com/ayush6624/go-chatgpt"
)

var (
	api_key = os.Getenv("FLASH_CARDS_OPEN_AI_API_KEY")
)

var client *chatgpt.Client

func init() {
	var err error
	client, err = chatgpt.NewClient(api_key)
	if err != nil {
		log.Fatal(err)
	}
}

func Request(req string) (string, error) {
	res, err := client.Send(context.Background(), &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT4,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleUser,
				Content: req,
			},
		},
	})
	if err != nil {
		return "", err
	}
	responce := ""
	if len(res.Choices) > 0 {
		responce = res.Choices[0].Message.Content
	}
	return responce, nil
}
