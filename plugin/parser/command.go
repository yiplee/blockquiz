package parser

import (
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
