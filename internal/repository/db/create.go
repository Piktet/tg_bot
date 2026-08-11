package db

import (
	"context"

	"github.com/Piktet/tg_bot/internal/model"
)

const q = `create table if not exists users (
    id bigint primary key,
    chat_id bigint,
    name text,
    created timestamp default now()
);

create table if not exists tasks (
    user_id bigint,
	chat_id bigint,
    task_id uuid,
    input_file_id uuid,
    output_file_id uuid,
    result text,
	result_short text,
    created timestamp default now()
);
`

// Create -  создание структуры таблиц
func Create(ctx context.Context, conn model.Connection) error {
	return conn.Execute(ctx, q)
}
