package speach

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

// SpeachConnection - структура подключения к speach.
type SpeachConnection struct {
	host    string
	RqUID   string
	AuthKey string

	token string
	err   error
}

// New - создание подключения к  speach.
func New(host string, rquid string, authkey string) *SpeachConnection {
	return &SpeachConnection{
		host:    host,
		RqUID:   rquid,
		AuthKey: authkey,
	}
}

// SpeachAuthResponse - возвращаемый токен.
type SpeachAuthResponse struct {
	Token   string `json:"access_token"`
	Expires int64  `json:"expires_at"`
}

// Connect - создает токен.
func (p *SpeachConnection) Connect(ctx context.Context) error {
	u, err := url.JoinPath(p.host, "/api/v2/oauth")
	if err != nil {
		logger.Log().Error("connection to speach - create path error", zap.Error(err))
		p.err = err
		return err
	}
	data := []byte("scope=SALUTE_SPEECH_PERS")
	r, err := http.NewRequestWithContext(ctx, "POST", u, bytes.NewReader(data))
	if err != nil {
		logger.Log().Error("connection to speach - create request", zap.Error(err))
		p.err = err
		return err
	}
	r.Header.Set("Authorization", "application/x-www-form-urlencoded")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("RqUID", p.RqUID)
	r.Header.Add("Authorization", "Bearer "+p.AuthKey)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		logger.Log().Error("connection to speach - send request", zap.Error(err))
		p.err = err
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("ошибка авторизации")
		logger.Log().Error("connection to speach - get response", zap.Error(err))
		p.err = err
		return err
	}

	var token SpeachAuthResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&token); err != nil {
		logger.Log().Error("connection to speach - decode result", zap.Error(err))
		p.err = err
		return err
	}

	p.token = token.Token

	t := time.Unix(token.Expires, 0)
	time.AfterFunc(time.Until(t), func() {
		if err := p.Connect(ctx); err != nil {
			logger.Log().Error("speach reconnect error", zap.Error(err))
			return
		}
	})

	return nil
}

// GetToken - возвращает токен, проверяет на ошибку получения токена.
func (p *SpeachConnection) GetToken() (string, error) {
	if p.err != nil {
		return "", p.err
	}
	return p.token, nil
}
