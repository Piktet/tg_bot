package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/Piktet/tg_bot/internal/logger"

	"go.uber.org/zap"
)

// chatConnection - структура подключения к chat.
type ChatConnection struct {
	host    string
	RqUID   string
	AuthKey string

	token string
	err   error
}

// New - создание подключения к  chat.
func New(host string, rquid string, authkey string) *ChatConnection {
	return &ChatConnection{
		host:    host,
		RqUID:   rquid,
		AuthKey: authkey,
	}
}

// chatAuthResponse - возвращаемый токен.
type chatAuthResponse struct {
	Token   string `json:"access_token"`
	Expires int64  `json:"expires_at"`
}

// Connect - создает токен.
func (p *ChatConnection) Connect(ctx context.Context) error {
	u, err := url.JoinPath(p.host, "/api/v2/oauth")
	if err != nil {
		logger.Log().Error("connection to chat - create path error", zap.Error(err))
		p.err = err
		return err
	}
	data := []byte("scope=GIGACHAT_API_PERS")
	r, err := http.NewRequestWithContext(ctx, "POST", u, bytes.NewReader(data))
	if err != nil {
		logger.Log().Error("connection to chat - create request", zap.Error(err))
		p.err = err
		return err
	}
	r.Header.Set("Authorization", "application/x-www-form-urlencoded")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("RqUID", p.RqUID)
	r.Header.Add("Authorization", "Basic "+p.AuthKey)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		logger.Log().Error("connection to chat - send request", zap.Error(err))
		p.err = err
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("ошибка авторизации")
		logger.Log().Error("connection to chat - get response", zap.Error(err))
		p.err = err
		return err
	}

	var token chatAuthResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&token); err != nil {
		logger.Log().Error("connection to chat - decode result", zap.Error(err))
		p.err = err
		return err
	}

	p.token = token.Token

	t := time.Unix(token.Expires, 0)
	time.AfterFunc(time.Until(t), func() {
		if err := p.Connect(ctx); err != nil {
			logger.Log().Error("chat reconnect error", zap.Error(err))
			return
		}
	})

	return nil
}

// GetToken - возвращает токен, проверяет на ошибку получения токена.
func (p *ChatConnection) GetToken() (string, error) {
	if p.err != nil {
		return "", p.err
	}
	return p.token, nil
}
