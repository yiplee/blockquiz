package parser

import (
	"github.com/yiplee/blockquiz/core"
)

func parseCommand(args Args) (*core.Command, error) {
	cmd := core.Command{
		Action: core.ActionUsage,
	}

	switch args.First() {
	case core.ActionSwitchChinese:
		cmd.Action = core.ActionSwitchChinese
	case core.ActionSwitchEnglish:
		cmd.Action = core.ActionSwitchEnglish
	case core.ActionShowCourse:
		cmd.Action = core.ActionShowCourse
	case core.ActionRandomCourse:
		cmd.Action = core.ActionRandomCourse
	case core.ActionShowQuestion:
		cmd.Action = core.ActionShowQuestion
	case core.ActionAnswerQuestion:
		/*
			> answer_question 1
			arg(1) present answer
		*/
		cmd.Action = core.ActionAnswerQuestion
		cmd.Answer, _ = args.GetInt(1)
	default:
		// a - z 算答题
		if runes := []byte(args.First()); len(runes) == 1 && len(args) == 1 {
			if r := runes[0]; r >= 'a' && r <= 'd' {
				cmd.Action = core.ActionAnswerQuestion
				cmd.Answer = int(r - 'a')
			}
		}
	}

	return &cmd, nil
}

func encodeCommand(cmd *core.Command) (args Args) {
	args = append(args, cmd.Action)

	if cmd.Action == core.ActionAnswerQuestion {
		args[0] = core.AnswerToString(cmd.Answer)
	}

	return
}
