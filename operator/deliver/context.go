package deliver

import (
	"context"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/store"
)

type commandContext struct {
	d              *Deliver
	cmd            *core.Command
	user           *core.User
	course         *core.Course
	question       *core.Question
	localizer      *localizer.Localizer
	conversationID string
	traceID        string
}

func (d *Deliver) prepareContext(ctx context.Context, cmd *core.Command) (*commandContext, error) {
	c := &commandContext{
		d:   d,
		cmd: cmd,
	}

	var err error
	c.user, err = d.users.FindMixinID(ctx, cmd.UserID)
	if store.IsErrNotFound(err) {
		c.user = &core.User{
			MixinID: cmd.UserID,
		}
		err = d.users.Create(ctx, c.user)
	}

	if err != nil {
		return nil, err
	}

	if cmd.Course > 0 {
		c.course, err = d.courses.Find(ctx, cmd.Course)
		if err != nil && !store.IsErrNotFound(err) {
			return nil, err
		}
	}

	if c.course != nil {
		c.question, _ = c.course.Question(cmd.Question)
	}

	c.localizer = localizer.WithLanguage(d.localizer, c.user.Language)
	c.conversationID = bot.UniqueConversationId(c.user.MixinID, d.config.ClientID)
	c.traceID = cmd.TraceID

	return c, nil
}
