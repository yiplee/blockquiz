package core

import (
	"context"
	"time"
)

const (
	ActionUsage          = "usage"           // 帮助
	ActionSwitchLanguage = "language"        // 切换语言
	ActionSwitchEnglish  = "en"              // 切换到英文
	ActionSwitchChinese  = "zh"              // 切换到中文
	ActionShowQuestion   = "show_question"   // 开始答题
	ActionAnswerQuestion = "answer_question" // 答题
)

const (
	CommandSourcePlainText = "plain_text"
	CommandSourceSnapshot  = "snapshot"
	CommandSourceAPI       = "api" // api
)

type (
	Command struct {
		ID        int64      `gorm:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time  `json:"created_at,omitempty"`
		UpdatedAt time.Time  `json:"updated_at,omitempty"`
		DeletedAt *time.Time `json:"deleted_at,omitempty"`
		TraceID   string     `gorm:"size:36" json:"trace_id,omitempty"`
		UserID    string     `gorm:"size:36" json:"user_id,omitempty"`
		Action    string     `gorm:"size:256" json:"action,omitempty"`
		Answer    int        `json:"answer,omitempty"`
		Source    string     `gorm:"size:64" json:"source,omitempty"`
	}

	CommandStore interface {
		Create(ctx context.Context, command *Command) error
		Delete(ctx context.Context, command *Command) error
		Deletes(ctx context.Context, commands []*Command) error
		ListPending(ctx context.Context, limit int) ([]*Command, error)
	}

	CommandParser interface {
		Parse(ctx context.Context, input string) ([]*Command, error)
		Encode(ctx context.Context, cmds ...*Command) string
	}
)
