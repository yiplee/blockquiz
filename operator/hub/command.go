package hub

import (
	"context"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/yiplee/blockquiz/core"
)

func (h *Hub) parseCommand(ctx context.Context, mixinID, traceID string, args Args) (*core.Command, error) {
	cmd := core.Command{
		TraceID: traceID,
		UserID:  mixinID,
	}

	switch args.First() {
	case core.ActionHelp, "?", "usage", "hi":
		cmd.Action = core.ActionHelp
	case core.ActionSwitchChinese:
		cmd.Action = core.ActionSwitchChinese
	case core.ActionSwitchEnglish:
		cmd.Action = core.ActionSwitchEnglish
	case core.ActionShowLesson:
		/*
			> show_lesson 1
			arg(1) present lesson id
		*/
		cmd.Action = core.ActionShowLesson
		cmd.Course, _ = args.GetInt64(1)
	case core.ActionShowQuestion:
		/*
			> show_question 1
			arg(1) present lesson id
		*/
		cmd.Action = core.ActionShowQuestion
		cmd.Course, _ = args.GetInt64(1)
	case core.ActionAnswerQuestion:
		/*
			> answer_question 1
			arg(1) present lesson id
		*/
		cmd.Action = core.ActionAnswerQuestion
		cmd.Course, _ = args.GetInt64(1)
	}

	return nil, nil
}

func parseLessonID(arg string) (int64, bool) {
	n, err := strconv.ParseInt(arg, 10, 64)
	if err != nil || n <= 0 {
		return 0, false
	}

	pflag.Arg(0)
	return n, false
}

func validateCommand(cmd *core.Command) error {
	return nil
}
