package db

import (
	"context"
	"unsafe"

	"github.com/Piktet/tg_bot/internal/logger"
	"github.com/Piktet/tg_bot/internal/model"

	"go.uber.org/zap"
)

const (
	qAddUser = "insert into users(id, chat_id, name, created) values ($1, $2, $3, current_timestamp) ON CONFLICT DO NOTHING"
	qAddTask = "insert into tasks(user_id, chat_id, task_id, input_file_id, output_file_id, result, result_short) values($1, $2, $3, $4, $5, $6, $7)"
)

// SaveUser - регистрация пользователя – запоминаем его идентификатор.
func AddUser(ctx context.Context, conn model.Connection, id, chatId int64, name string) error {
	if err := conn.Execute(ctx, qAddUser, id, chatId, name); err != nil {
		logger.Log().Error("error AddUser - Query", zap.Error(err))
		return err
	}
	return nil
}

func AddTask(ctx context.Context, conn model.Connection, x *model.SpeachTaskResponse, short string) error {
	return conn.Execute(ctx, qAddTask, x.User, x.ChatID, x.TaskID, x.InFileID, x.OutFileID, unsafe.String(unsafe.SliceData(x.Output), len(x.Output)), short)
}
