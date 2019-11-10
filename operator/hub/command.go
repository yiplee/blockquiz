package hub

import (
	"context"
	"strconv"

	"github.com/yiplee/blockquiz/core"
)

func (h *Hub) handleCommand(ctx context.Context, mixinID, traceID string, args []string) error {
	cmd := core.Command{
		TraceID: traceID,
		UserID:  mixinID,
	}

	switch args[0] {
	case core.ActionHelp, "?", "usage", "hi":
		cmd.Action = core.ActionHelp
	case core.ActionSwitchChinese:
		cmd.Action = core.ActionSwitchChinese
	case core.ActionSwitchEnglish:
		cmd.Action = core.ActionSwitchEnglish
	case core.ActionShowLesson:
		/*
			> show_lesson 1
			arg(0) present lesson id
		*/
		cmd.Action = core.ActionShowLesson

	}

	return nil
}

func parseLessonID(arg string) (int64, bool) {
	n, err := strconv.ParseInt(arg, 10, 64)
	if err != nil || n <= 0 {
		return 0, false
	}

	return n, false
}
