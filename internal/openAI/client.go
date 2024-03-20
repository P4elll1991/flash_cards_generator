package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ayush6624/go-chatgpt"
	chatgpt_errors "github.com/ayush6624/go-chatgpt/utils"
)

var (
	api_key       = os.Getenv("FLASH_CARDS_OPEN_AI_API_KEY")
	chatgpt_model = os.Getenv("FLASH_CARDS_OPEN_AI_CHATGPT_MODEL")
)

var client *Client

func init() {
	var err error
	fmt.Println("chatgpt_model", chatgpt_model)
	client, err = NewClient(api_key)
	if err != nil {
		log.Fatal(err)
	}
}

func Request(req string) (string, error) {
	res, err := client.Send(context.Background(), &chatgpt.ChatCompletionRequest{
		Model: chatgpt.ChatGPTModel(chatgpt_model),
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

const (
	apiURL = "https://api.openai.com/v1"
)

type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Config
	config *Config
}

type Config struct {
	// Base URL for API requests.
	BaseURL string

	// API Key (Required)
	APIKey string

	// Organization ID (Optional)
	OrganizationID string
}

func NewClient(apikey string) (*Client, error) {
	if apikey == "" {
		return nil, chatgpt_errors.ErrAPIKeyRequired
	}

	return &Client{
		client: &http.Client{},
		config: &Config{
			BaseURL: apiURL,
			APIKey:  apikey,
		},
	}, nil
}

func NewClientWithConfig(config *Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, chatgpt_errors.ErrAPIKeyRequired
	}

	return &Client{
		client: &http.Client{},
		config: config,
	}, nil
}

func (c *Client) sendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	if c.config.OrganizationID != "" {
		req.Header.Set("OpenAI-Organization", c.config.OrganizationID)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		// Parse body
		var errMessage interface{}
		if err := json.NewDecoder(res.Body).Decode(&errMessage); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("api request failed: status Code: %d %s %s Message: %+v", res.StatusCode, res.Status, res.Request.URL, errMessage)
	}

	return res, nil
}

func (c *Client) Send(ctx context.Context, req *chatgpt.ChatCompletionRequest) (*chatgpt.ChatResponse, error) {
	reqBytes, _ := json.Marshal(req)

	endpoint := "/chat/completions"
	httpReq, err := http.NewRequest("POST", c.config.BaseURL+endpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	httpReq = httpReq.WithContext(ctx)

	res, err := c.sendRequest(ctx, httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var chatResponse chatgpt.ChatResponse
	if err := json.NewDecoder(res.Body).Decode(&chatResponse); err != nil {
		return nil, err
	}

	return &chatResponse, nil
}
