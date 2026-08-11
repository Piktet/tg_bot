package bot

import (
	"context"
	"time"

	"github.com/Piktet/tg_bot/internal/model"
	"github.com/Piktet/tg_bot/internal/repository/db"

	tele "gopkg.in/telebot.v3"
)

type SpeachTask interface {
	AddTask(*model.SpeachTaskData)
	GetTaskResponse() chan *model.SpeachTaskResponse
}

type chatTask interface {
	GetShort(context.Context, []byte) (string, error)
	GetChat(context.Context, string) (string, error)
}

type DBTask interface {
	AddUser(context.Context, int64, int64, string) error
}

type Bot struct {
	b    *tele.Bot
	user *tele.User

	speachTaskProcessor SpeachTask
	chatProcessor       chatTask
	conn                model.Connection

	chCommand chan *Command
}

func New(token string, conn model.Connection, s SpeachTask, c chatTask) *Bot {
	b, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil
	}
	x := &Bot{
		b:                   b,
		speachTaskProcessor: s,
		chatProcessor:       c,
		conn:                conn,
	}

	b.Handle("/start", x.HandlerStart)
	b.Handle("/list", x.HandlerList)
	b.Handle("/get", x.HandlerGet)
	b.Handle("/find", x.HandlerFind)
	b.Handle("/chat", x.Handlerchat)

	b.Handle(tele.OnAudio, x.HandlerOnAudio)
	b.Handle(tele.OnAudio, x.HandlerOnVoice)
	b.Handle(tele.OnAudio, x.HandlerOnText)

	return x
}

func (p *Bot) getTaskResultProcess(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case task := <-p.speachTaskProcessor.GetTaskResponse():
			if short, err := p.chatProcessor.GetShort(ctx, task.Output); err == nil {
				db.AddTask(ctx, p.conn, task, short)
			}
		}
	}
}
