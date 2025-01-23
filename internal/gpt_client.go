package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type GPTRequest struct {
	Messages []map[string]string `json:"messages"`
	Model    string              `json:"model"`
}

type GPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func ReviewCode(filePath string) (string, error) {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("⚠️  Warning: Could not load .env file")
	}

	// Read the file to be reviewed
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("❌ Could not read file %s: %w", filePath, err)
	}

	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", errors.New("❌ OPENAI_API_KEY environment variable is not set")
	}

	// Prepare request payload
	payload := GPTRequest{
		Model: "gpt-4",
		Messages: []map[string]string{
			{"role": "system", "content": "너는 전문 코드 리뷰어야. 코드 리뷰를 해주고 한국어로 답변해줘."},
			{"role": "user", "content": fmt.Sprintf("다음 코드를 리뷰해줘:\n\n%s", string(content))},
		},
	}

	body, _ := json.Marshal(payload)

	// Create HTTP request
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP responses
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("❌ API call failed with status %s: %s", resp.Status, string(responseBody))
	}

	// Parse response body
	responseBody, _ := ioutil.ReadAll(resp.Body)
	var gptResponse GPTResponse
	err = json.Unmarshal(responseBody, &gptResponse)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse API response: %w", err)
	}

	// Return GPT response content
	if len(gptResponse.Choices) == 0 {
		return "", errors.New("❌ No choices returned in API response")
	}

	return fmt.Sprintf("✅ Code Review Result:\n%s", gptResponse.Choices[0].Message.Content), nil
}
