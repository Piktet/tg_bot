package bot

import (
	"bytes"
	"context"
	"unsafe"

	"github.com/Piktet/tg_bot/internal/logger"
	"github.com/Piktet/tg_bot/internal/model"
	"github.com/Piktet/tg_bot/internal/repository/db"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	tele "gopkg.in/telebot.v3"
)

type Command struct {
	Name string
	fn   func(context.Context) error
}

func (p *Bot) Start(ctx context.Context, cnt, size int) error {
	p.chCommand = make(chan *Command, size)
	defer close(p.chCommand)
	var wg errgroup.Group
	for range cnt {
		wg.Go(func() error {
			return p.worker(ctx)
		})
	}
	wg.Go(func() error {
		go p.b.Start()
		<-ctx.Done()
		p.b.Stop()
		return nil
	})

	wg.Go(func() error {
		return p.getTaskResultProcess(ctx)
	})
	return wg.Wait()

}

func (p *Bot) addCommand(c tele.Context, cmd *Command) error {
	select {
	case p.chCommand <- cmd:
		logger.Log().Info("command added", zap.String("name", cmd.Name))
		return nil
	default:
		return c.Send("try later")
	}
}

func (p *Bot) worker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case cmd := <-p.chCommand:
			cmd.fn(ctx)
		}
	}
}

// Start - регистрация пользователя – запоминаем его идентификатор.
func (p *Bot) HandlerStart(c tele.Context) error {
	return p.addCommand(c, &Command{fn: func(ctx context.Context) error {
		return db.AddUser(ctx, p.conn, c.Sender().ID, c.Chat().ID, c.Sender().Username)
	}})
}

// Get - получение текста встречи.
func (p *Bot) HandlerGet(c tele.Context) error {

	args := c.Args()
	if len(args) < 1 {
		return c.Send("/get <id>")
	}

	id := args[0]

	return p.addCommand(c, &Command{fn: func(ctx context.Context) error {

		user := c.Sender().ID
		result, err := db.GetUserFileItem(ctx, p.conn, user, id)
		if err != nil {
			return err
		}

		c.Send(result)

		return err
	}})
}

// List - список сохраненных встреч.
func (p *Bot) HandlerList(c tele.Context) error {
	return p.addCommand(c, &Command{fn: func(ctx context.Context) error {
		result, err := db.GetUserFile(ctx, p.conn, c.Sender().ID)
		if err != nil {
			return err
		}

		c.Send(result)

		return err
	}})
}

// Find - поиск встречи по ключевым словам.
func (p *Bot) HandlerFind(c tele.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return c.Send("/find <word>")
	}

	word := args[0]
	return p.addCommand(c, &Command{fn: func(ctx context.Context) error {

		list, err := db.GetFileByWord(ctx, p.conn, c.Sender().ID, word)

		c.Send(list)
		return err
	}})
}

// chat - запрос к GigaChat.
func (p *Bot) Handlerchat(c tele.Context) error {
	return p.addCommand(c, &Command{fn: func(ctx context.Context) error {

		x, err := p.chatProcessor.GetChat(ctx, c.Message().Text)
		if err != nil {
			return err
		}
		c.Send(x)

		return nil
	}})
}

func (p *Bot) HandlerOnVoice(c tele.Context) error {
	return p.addCommand(c, &Command{fn: func(ctx context.Context) error {
		p.speachTaskProcessor.AddTask(&model.SpeachTaskData{
			User:   c.Sender().ID,
			ChatID: c.Chat().ID,
			Input:  c.Message().Voice.FileReader})
		return nil
	}})
}

func (p *Bot) HandlerOnAudio(c tele.Context) error {
	return p.addCommand(c, &Command{fn: func(ctx context.Context) error {
		p.speachTaskProcessor.AddTask(&model.SpeachTaskData{
			User:   c.Sender().ID,
			ChatID: c.Chat().ID,
			Input:  c.Message().Audio.FileReader})
		return nil
	}})
}

func (p *Bot) HandlerOnText(c tele.Context) error {
	return p.addCommand(c, &Command{
		fn: func(ctx context.Context) error {
			p.speachTaskProcessor.AddTask(&model.SpeachTaskData{
				User:   c.Sender().ID,
				ChatID: c.Chat().ID,
				Input:  bytes.NewReader(unsafe.Slice(unsafe.StringData(c.Message().Text), len(c.Message().Text)))})
			return nil
		}})
}
