package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", errors.New("OPENAI_API_KEY 환경 변수가 설정되지 않았습니다")
	}

	payload := GPTRequest{
		Model: "gpt-4",
		Messages: []map[string]string{
			{"role": "system", "content": "You are a professional code reviewer."},
			{"role": "user", "content": fmt.Sprintf("Review the following code:\n\n%s", string(content))},
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 호출 실패: %s", resp.Status)
	}

	responseBody, _ := ioutil.ReadAll(resp.Body)
	var gptResponse GPTResponse
	err = json.Unmarshal(responseBody, &gptResponse)
	if err != nil {
		return "", err
	}

	return gptResponse.Choices[0].Message.Content, nil
}
