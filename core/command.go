package core

import (
	"context"
	"time"
)

const (
	ActionUsage          = "usage"           // 帮助
	ActionSwitchEnglish  = "en"              // 切换到英文
	ActionSwitchChinese  = "zh"              // 切换到中文
	ActionShowCourse     = "show_course"     // 开始课程
	ActionRandomCourse   = "random_course"   // 随机课程
	ActionShowQuestion   = "show_question"   // 开始答题
	ActionAnswerQuestion = "answer_question" // 答题
	ActionRequestCoin    = "coin"            // 请求答题币
)

type (
	Command struct {
		TraceID   string    `gorm:"size:36" json:"trace_id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UserID    string    `gorm:"size:36" json:"user_id,omitempty"`
		Action    string    `gorm:"size:256" json:"action,omitempty"`
		Course    int64     `json:"course,omitempty"`
		Question  int       `json:"question_number,omitempty"`
		Answer    int       `json:"answer,omitempty"`
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
