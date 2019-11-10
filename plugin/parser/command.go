package parser

import (
	"fmt"
	"strconv"

	"github.com/yiplee/blockquiz/core"
)

func parseCommand(args Args) (*core.Command, error) {
	cmd := core.Command{}

	switch args.First() {
	case core.ActionUsage, "?", "hi":
		cmd.Action = core.ActionUsage
	case core.ActionSwitchChinese:
		cmd.Action = core.ActionSwitchChinese
	case core.ActionSwitchEnglish:
		cmd.Action = core.ActionSwitchEnglish
	case core.ActionShowCourse:
		/*
			> show_course 65
			arg(1) present course id
		*/
		cmd.Action = core.ActionShowCourse
		cmd.Course, _ = args.GetInt64(1)
	case core.ActionRandomCourse:
		cmd.Action = core.ActionRandomCourse
	case core.ActionShowQuestion:
		/*
			> show_question 65 2
			arg(1) present course id
			arg(2) present question serial number
		*/
		cmd.Action = core.ActionShowQuestion
		cmd.Course, _ = args.GetInt64(1)
		cmd.Question, _ = args.GetInt(2)
	case core.ActionAnswerQuestion:
		/*
			> answer_question 65 2 1
			arg(1) present course id
			arg(2) present question serial number
			arg(3) present answer
		*/
		cmd.Action = core.ActionAnswerQuestion
		cmd.Course, _ = args.GetInt64(1)
		cmd.Question, _ = args.GetInt(2)
		cmd.Answer, _ = args.GetInt(3)
	case core.ActionRequestCoin:
		cmd.Action = core.ActionRequestCoin
	default:
		return nil, fmt.Errorf("unknown action %s", args.First())
	}

	return &cmd, nil
}

func encodeCommand(cmd *core.Command) (args Args) {
	args = append(args, cmd.Action)

	switch cmd.Action {
	case core.ActionShowCourse:
		args = append(args, strconv.FormatInt(cmd.Course, 10))
	case core.ActionShowQuestion:
		args = append(args, strconv.FormatInt(cmd.Course, 10))
		args = append(args, strconv.Itoa(cmd.Question))
	case core.ActionAnswerQuestion:
		args = append(args, strconv.FormatInt(cmd.Course, 10))
		args = append(args, strconv.Itoa(cmd.Question))
		args = append(args, strconv.Itoa(cmd.Answer))
	}

	return
}
