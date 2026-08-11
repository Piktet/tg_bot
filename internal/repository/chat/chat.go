package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"unsafe"

	"github.com/Piktet/tg_bot/internal/logger"

	"go.uber.org/zap"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

const prompt = "сделай краткую выжимку из текста %s"

func GetShort(ctx context.Context, host, token string, text []byte) (string, error) {
	return getData(ctx, host, token, "user", unsafe.String(unsafe.SliceData(text), len(text)))
}

func GetChat(ctx context.Context, host, token string, text string) (string, error) {
	return getData(ctx, host, token, "assistant", text)
}

// Voice - загрузить файл для отправки в SaluteSpeech API.
func getData(ctx context.Context, host, token, role string, text string) (string, error) {
	u, err := url.JoinPath(host, "/api/v1/chat/completions")
	if err != nil {
		logger.Log().Error("UploadVoice - create path error", zap.Error(err))
		return "", err
	}

	data := &ChatRequest{
		Model: "GigaChat",
		Messages: []*Message{{
			Role:    role,
			Content: fmt.Sprintf(prompt, text),
		}},
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	r, err := http.NewRequestWithContext(ctx, "POST", u, bytes.NewBuffer(body))
	if err != nil {
		logger.Log().Error("UploadVoice - create request", zap.Error(err))
		return "", err
	}

	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		logger.Log().Error("UploadVoice - send request", zap.Error(err))
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("ошибка авторизации")
		logger.Log().Error("UploadVoice - get response", zap.Error(err))
		return "", err
	}

	var x ChatResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&x); err != nil {
		logger.Log().Error("UploadVoice - decode result", zap.Error(err))
		return "", err
	}

	if len(x.Choices) > 0 {
		return x.Choices[0].Message.Content, nil
	}

	return "", nil

}
