package model

import (
	"context"
	"database/sql"
	"io"
	"time"
)

type Meet struct {
	name    string
	created time.Time
	keys    []string
}

type ResultStatusType string

const (
	SpeachResultStatusNew      ResultStatusType = "NEW"
	SpeachResultStatusRunning  ResultStatusType = "RUNNING"
	SpeachResultStatusCanceled ResultStatusType = "CANCELED"
	SpeachResultStatusDone     ResultStatusType = "DONE"
	SpeachResultStatusError    ResultStatusType = "ERROR"
	SpeachResultStatusEmpty    ResultStatusType = ""
)

type SpeachUploadResponse struct {
	Status int `json:"status"`
	Result struct {
		FileID string `json:"request_file_id"`
	} `json:"result"`
}

type SpeachCreateTaskOptionRequest struct {
	AudioEncoding string `json:"audio_encoding"`
}
type SpeachCreateTaskRequest struct {
	FileID  string                        `json:"request_file_id"`
	Options SpeachCreateTaskOptionRequest `json:"options"`
}
type SpeachCreateTaskResultResponse struct {
	ID      string           `json:"id"`
	Created string           `json:"created_at"`
	Updated string           `json:"updated_at"`
	Status  ResultStatusType `json:"status"`
	FileID  string           `json:"response_file_id,omitempty"`
	Error   string           `json:"error,omitempty"`
}
type SpeachCreateTaskResponse struct {
	Status int                            `json:"status"`
	Result SpeachCreateTaskResultResponse `json:"result"`
}

// SpeachTaskData - данные для задачи в SaluteSpech
type SpeachTaskData struct {
	User   int64
	ChatID int64
	Input  io.Reader
}

// SpeachTaskData - данные для задачи в SaluteSpech
type SpeachTaskResponse struct {
	*SpeachTaskData
	InFileID  string
	OutFileID string
	TaskID    string
	Output    []byte
}

type FileInfo struct {
	ChatID string
	TaskID string
	FileID string
	Data   string
}

type Words struct {
	ID, Word string
}

// Connection - интерфейс подключения к БД
type Connection interface {
	Query(context.Context, string, ...any) (*sql.Rows, error)
	Execute(context.Context, string, ...any) error
	BeginTx(context.Context) (*sql.Tx, error)
}
