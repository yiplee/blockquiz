package parser

import (
	"strconv"

	"github.com/yiplee/blockquiz/core"
)

func parseCommand(args Args) (*core.Command, error) {
	cmd := core.Command{
		Action: core.ActionUsage,
	}

	switch args.First() {
	case core.ActionSwitchLanguage, "切换语言":
		cmd.Action = core.ActionSwitchLanguage
	case core.ActionSwitchChinese, "中文":
		cmd.Action = core.ActionSwitchChinese
	case core.ActionSwitchEnglish, "english":
		cmd.Action = core.ActionSwitchEnglish
	case core.ActionShowQuestion, "答题", "开始答题":
		cmd.Action = core.ActionShowQuestion
	default:
		if question, ok := args.GetInt(0); ok {
			if answer, ok := args.Get(1); ok {
				if runes := []byte(answer); len(runes) == 1 {
					if r := runes[0]; r >= 'a' && r <= 'd' {
						cmd.Action = core.ActionAnswerQuestion
						cmd.Answer = int(r - 'a')
						cmd.Question = question - 1
					}
				}
			}
		}
	}

	return &cmd, nil
}

func encodeCommand(cmd *core.Command) Args {
	switch cmd.Action {
	case core.ActionAnswerQuestion:
		return Args{
			strconv.Itoa(cmd.Question + 1),
			core.AnswerToString(cmd.Answer),
		}
	default:
		return Args{cmd.Action}
	}
}
