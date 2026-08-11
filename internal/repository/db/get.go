package db

import (
	"context"
	"errors"

	"github.com/Piktet/tg_bot/internal/logger"
	"github.com/Piktet/tg_bot/internal/model"

	"go.uber.org/zap"
)

const (
	qGetUserFileList = "select chat_id, task_id, input_file_id from tasks where user_id = $1 order by created"
	qGetUserFileItem = "select chat_id, task_id, input_file_id, result from tasks where user_id = $1 and file_id = $2"
	qGetFileWord     = "select chat_id, task_id, input_file_id from tasks where user_id = $1 and $2 ~ result or user_id = $3 and $4 ~ result_short"
)

func GetUserFile(ctx context.Context, conn model.Connection, id int64) ([]*model.FileInfo, error) {

	rows, err := conn.Query(ctx, qGetUserFileList, id)
	if err != nil {
		logger.Log().Error("error GetUserFile - Query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	x := make([]*model.FileInfo, 0)

	for rows.Next() {
		var item model.FileInfo
		if err := rows.Scan(&item.ChatID, &item.TaskID, &item.FileID); err != nil {
			logger.Log().Error("error GetUserFile - Scan", zap.Error(err))
			return nil, err
		}
		x = append(x, &item)
	}
	if rows.Err() != nil {
		logger.Log().Error("error GetUserFile - Finish", zap.Error(err))
		return nil, err
	}

	return x, nil
}

func GetUserFileItem(ctx context.Context, conn model.Connection, id int64, file string) (*model.FileInfo, error) {

	rows, err := conn.Query(ctx, qGetUserFileItem, id, file)
	if err != nil {
		logger.Log().Error("error GetUserFile - Query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var x model.FileInfo
	if !rows.Next() {
		return nil, errors.New("not found")
	}
	if err := rows.Scan(&x.ChatID, &x.TaskID, &x.FileID, &x.Data); err != nil {
		logger.Log().Error("error GetUserFile - Scan", zap.Error(err))
		return nil, err
	}

	if rows.Err() != nil {
		logger.Log().Error("error GetUserFile - Finish", zap.Error(err))
		return nil, err
	}

	return &x, nil
}

func GetFileByWord(ctx context.Context, conn model.Connection, id int64, word string) ([]*model.FileInfo, error) {

	rows, err := conn.Query(ctx, qGetFileWord, id, word, id, word)
	if err != nil {
		logger.Log().Error("error GetUserFile - Query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	x := make([]*model.FileInfo, 0)

	for rows.Next() {
		var item model.FileInfo
		if err := rows.Scan(&item.ChatID, &item.TaskID, &item.FileID); err != nil {
			logger.Log().Error("error GetUserFile - Scan", zap.Error(err))
			return nil, err
		}
		x = append(x, &item)
	}
	if rows.Err() != nil {
		logger.Log().Error("error GetUserFile - Finish", zap.Error(err))
		return nil, err
	}

	return x, nil
}
