package core

import (
	"context"
	"time"
)

const (
	ActionHelp           = "help"            // 帮助
	ActionSwitchEnglish  = "en"              // 切换到英文
	ActionSwitchChinese  = "zh"              // 切换到中文
	ActionSubscribe      = "subscribe"       // 新用户
	ActionShowLesson     = "show_lesson"     // 开始课程
	ActionShowQuestion   = "show_question"   // 开始答题
	ActionAnswerQuestion = "answer_question" // 答题
	ActionNextQuestion   = "next_question"   // 下一题
	ActionNextLesson     = "next_lesson"     // 下一课
	ActionRequestCoin    = "coin"            // 请求答题币
)

type Command struct {
	ID        int64      `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	DeletedAt *time.Time `gorm:"INDEX" json:"deleted_at,omitempty"`
	TraceID   string     `gorm:"size:36" json:"trace_id,omitempty"`
	UserID    string     `gorm:"size:36" json:"user_id,omitempty"`
	Action    string     `gorm:"size:256" json:"action,omitempty"`
	Chapter   int64      `json:"chapter,omitempty"`
	Question  int64      `json:"question_number,omitempty"`
	Answer    Answer     `json:"answer,omitempty"`
}

type CommandStore interface {
	Create(ctx context.Context, command *Command) error
	Delete(ctx context.Context, id int64) error
	Deletes(ctx context.Context, ids []int64) error
	List(ctx context.Context, limit int) ([]*Command, error)
}
