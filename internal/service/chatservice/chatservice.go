package chatservice

import (
	"context"
	"sync"

	"github.com/Piktet/tg_bot/internal/repository/chat"
)

type ChatService struct {
	mx   sync.Mutex
	host string
	conn *chat.ChatConnection
}

func New(host string, conn *chat.ChatConnection) *ChatService {
	return &ChatService{
		host: host,
		conn: conn,
	}
}

func (p *ChatService) GetShort(ctx context.Context, full []byte) (string, error) {
	p.mx.Lock()
	defer p.mx.Unlock()
	token, err := p.conn.GetToken()
	if err != nil {
		return "", err
	}

	return chat.GetShort(ctx, p.host, token, full)

}

func (p *ChatService) GetChat(ctx context.Context, text string) (string, error) {
	p.mx.Lock()
	defer p.mx.Unlock()
	token, err := p.conn.GetToken()
	if err != nil {
		return "", err
	}

	return chat.GetChat(ctx, p.host, token, text)

}
