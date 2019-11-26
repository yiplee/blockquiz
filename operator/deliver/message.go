package deliver

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/fox-one/pkg/text/localizer"
	"github.com/fox-one/pkg/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/thirdparty/bot-api-go-client"
)

// func (c *commandContext) paymentButtonAction(ctx context.Context, traceID string, cmds ...*core.Command) string {
// 	uri, _ := url.Parse("mixin://pay")
// 	query := uri.Query()
// 	query.Set("recipient", c.d.config.ClientID)
// 	query.Set("asset", c.d.config.CoinAsset)
// 	query.Set("amount", c.d.config.CoinAmount.Truncate(8).String())
// 	query.Set("trace", traceID)
// 	memo := c.d.parser.Encode(ctx, cmds...)
// 	query.Set("memo", memo)
// 	uri.RawQuery = query.Encode()
// 	return uri.String()
// }

func (c *commandContext) inputButtonAction(ctx context.Context, cmds ...*core.Command) string {
	return fmt.Sprintf("input:%s", c.d.parser.Encode(ctx, cmds...))
}

type button struct {
	Label  string `json:"label,omitempty"`
	Color  string `json:"color,omitempty"`
	Action string `json:"action,omitempty"`
}

func (c *commandContext) newButton(label, action string) button {
	return button{
		Label:  label,
		Color:  c.d.config.ButtonColor,
		Action: action,
	}
}

/*
新用户（还没设置语言）来的时候
发送一组设置语言的按钮，点击之后会发送 usage 信息和一个随机课程给用户
*/
func (c *commandContext) selectLanguage(ctx context.Context, next *core.Command) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "APP_BUTTON_GROUP",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "select language"),
	}

	var buttons []button
	for _, lang := range []string{core.ActionSwitchEnglish, core.ActionSwitchChinese} {
		cmds := []*core.Command{
			{Action: lang},
		}

		if next != nil {
			cmds = append(cmds, next)
		}

		l := localizer.WithLanguage(c.Localizer(), lang)

		buttons = append(buttons, c.newButton(
			l.MustLocalize("select_language"),
			c.inputButtonAction(ctx, cmds...),
		))
	}

	data, _ := jsoniter.Marshal(buttons)
	req.Data = base64.StdEncoding.EncodeToString(data)
	return req
}

func (c *commandContext) languageSwitched(ctx context.Context) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "PLAIN_TEXT",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "language switched"),
	}

	data := c.Localizer().MustLocalize("language_switched")
	req.Data = base64.StdEncoding.EncodeToString([]byte(data))
	return req
}

func (c *commandContext) showUsage(ctx context.Context, finish bool) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "PLAIN_TEXT",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "show usage"),
	}

	id := "usage_start_task"
	if finish {
		id = "usage_finish_task"
	}

	data := c.Localizer().MustLocalize(id)
	req.Data = base64.StdEncoding.EncodeToString([]byte(data))
	return req
}

func (c *commandContext) showUsageButtons(ctx context.Context, finish bool) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "APP_BUTTON_GROUP",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "show usage buttons"),
	}

	buttons := []button{
		c.newButton(
			c.Localizer().MustLocalize("switch_language"),
			c.inputButtonAction(ctx, &core.Command{
				Action: core.ActionSwitchLanguage,
			}),
		),
	}

	if !finish {
		buttons = append(buttons, c.newButton(
			c.Localizer().MustLocalize("show_question"),
			c.inputButtonAction(ctx, &core.Command{
				Action: core.ActionShowQuestion,
			}),
		))
	}

	data, _ := jsoniter.Marshal(buttons)
	req.Data = base64.StdEncoding.EncodeToString(data)
	return req
}

func (c *commandContext) showMissingCourse(ctx context.Context) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "PLAIN_TEXT",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "missing course"),
	}

	data := c.Localizer().MustLocalize("cannot_find_course")
	req.Data = base64.StdEncoding.EncodeToString([]byte(data))
	return req
}

func (c *commandContext) showCourseContent(ctx context.Context) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "PLAIN_TEXT",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "show course content"),
	}

	course := c.course

	var buf bytes.Buffer
	if c.task.Info != "" {
		fmt.Fprintln(&buf, c.task.Info)
		fmt.Fprintln(&buf) // 换行
	}

	if c.course.Title != "" {
		fmt.Fprintln(&buf, c.course.Title)
		fmt.Fprintln(&buf) // 换行
	}

	if course.URL == "" {
		fmt.Fprintln(&buf, course.Content)
	} else {
		fmt.Fprintln(&buf, course.Summary)
	}

	req.Data = base64.StdEncoding.EncodeToString(buf.Bytes())
	return req
}

func (c *commandContext) showCourseButtons(ctx context.Context) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "APP_BUTTON_GROUP",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "show course buttons"),
	}

	var buttons []button
	course := c.course
	if course.URL != "" {
		buttons = append(buttons, c.newButton(
			c.Localizer().MustLocalize("show_course"),
			course.URL,
		))
	}

	showQuestion := &core.Command{
		Action: core.ActionShowQuestion,
	}

	buttons = append(buttons, c.newButton(
		c.Localizer().MustLocalize("show_question"),
		c.inputButtonAction(ctx, showQuestion),
	))

	data, _ := jsoniter.Marshal(buttons)
	req.Data = base64.StdEncoding.EncodeToString(data)
	return req
}

func (c *commandContext) showQuestionContent(ctx context.Context) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "PLAIN_TEXT",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "show question content"),
	}

	task := c.task

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d/%d ", task.Question+1, len(c.course.Questions))
	fmt.Fprintln(&buf, c.question.Content)
	fmt.Fprintln(&buf)
	for idx, choice := range c.question.Choices {
		fmt.Fprintf(&buf, "%s %s\n", core.AnswerToString(idx), choice)
	}

	req.Data = base64.StdEncoding.EncodeToString(buf.Bytes())
	return req
}

func (c *commandContext) showQuestionChoiceButtons(ctx context.Context) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "APP_BUTTON_GROUP",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "show question buttons"),
	}

	buttons := make([]button, len(c.question.Choices))
	for idx := range buttons {
		cmd := &core.Command{
			Action:   core.ActionAnswerQuestion,
			Answer:   idx,
			Question: c.task.Question,
		}
		buttons[idx] = c.newButton(
			core.AnswerToString(idx),
			c.inputButtonAction(ctx, cmd),
		)
	}

	data, _ := jsoniter.Marshal(buttons)
	req.Data = base64.StdEncoding.EncodeToString(data)
	return req
}

func (c *commandContext) showAnswerFeedback(ctx context.Context, right bool) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "PLAIN_TEXT",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "answer feedback"),
	}

	var data string

	if right {
		data = c.Localizer().MustLocalize("answer_right")
	} else {
		if blocked, dur := c.task.IsBlocked(); blocked {
			minutes := int(dur.Minutes())
			if minutes == 0 {
				minutes = 1
			}
			data = c.Localizer().MustLocalize("answer_wrong_with_wait", "wait", strconv.Itoa(minutes))
		} else {
			data = c.Localizer().MustLocalize("answer_wrong")
		}
	}
	req.Data = base64.StdEncoding.EncodeToString([]byte(data))
	return req
}

func (c *commandContext) showWaitBlock(ctx context.Context) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "PLAIN_TEXT",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "wait block"),
	}

	_, dur := c.task.IsBlocked()
	minutes := int(dur.Minutes())
	if minutes == 0 {
		minutes = 1
	}
	data := c.Localizer().MustLocalize("wait_block", "wait", strconv.Itoa(minutes))
	req.Data = base64.StdEncoding.EncodeToString([]byte(data))
	return req
}

func (c *commandContext) showFinishCourse(ctx context.Context) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "PLAIN_TEXT",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "finish course"),
	}

	data := c.Localizer().MustLocalize("usage_finish_task")
	req.Data = base64.StdEncoding.EncodeToString([]byte(data))
	return req
}

func (c *commandContext) showNextQuestionButton(ctx context.Context, nextQuestion int) *bot.MessageRequest {
	req := &bot.MessageRequest{
		Category:       "APP_BUTTON_GROUP",
		RecipientId:    c.user.MixinID,
		ConversationId: c.conversationID,
		MessageId:      uuid.Modify(c.traceID, "next question button"),
	}

	cmd := &core.Command{
		Action: core.ActionShowQuestion,
	}

	buttons := []button{c.newButton(
		c.Localizer().MustLocalize("next_question"),
		c.inputButtonAction(ctx, cmd),
	)}

	data, _ := jsoniter.Marshal(buttons)
	req.Data = base64.StdEncoding.EncodeToString(data)
	return req
}
