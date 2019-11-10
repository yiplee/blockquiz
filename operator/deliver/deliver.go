package deliver

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/fox-one/pkg/uuid"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/store"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const limit = 100

type Deliver struct {
	users     core.UserStore
	commands  core.CommandStore
	parser    core.CommandParser
	courses   core.CourseStore
	wallets   core.WalletStore
	localizer *localizer.Localizer
	config    Config
}

func New(
	users core.UserStore,
	commands core.CommandStore,
	parser core.CommandParser,
	courses core.CourseStore,
	wallets core.WalletStore,
	localizer *localizer.Localizer,
	config Config,
) *Deliver {
	if _, err := govalidator.ValidateStruct(config); err != nil {
		panic(err)
	}

	return &Deliver{
		users:     users,
		commands:  commands,
		parser:    parser,
		courses:   courses,
		wallets:   wallets,
		localizer: localizer,
		config:    config,
	}
}

func (d *Deliver) Run(ctx context.Context, dur time.Duration) error {
	log := logger.FromContext(ctx).WithField("operator", "deliver")
	ctx = logger.WithContext(ctx, log)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			_ = d.run(ctx)
		}
	}
}

func (d *Deliver) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	list, err := d.commands.ListPending(ctx, limit)
	if err != nil {
		log.WithError(err).Error("list pending commands")
		return err
	}

	// group by userID
	group := make(map[string][]*core.Command)
	for _, cmd := range list {
		group[cmd.UserID] = append(group[cmd.UserID], cmd)
	}

	var g errgroup.Group
	var sem = semaphore.NewWeighted(3)

	for userID, cmds := range group {
		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}

		userID, cmds := userID, cmds
		g.Go(func() error {
			defer sem.Release(1)
			return d.post(ctx, userID, cmds)
		})
	}

	return g.Wait()
}

func (d *Deliver) post(ctx context.Context, userID string, cmds []*core.Command) error {
	log := logger.FromContext(ctx).WithField("user_id", userID)

	var requests []*bot.MessageRequest

	for _, cmd := range cmds {
		reqs, err := d.handleCommand(ctx, cmd)
		if err != nil {
			log.WithError(err).Error("handle command")
			return err
		}

		requests = append(requests, reqs...)
	}

	if len(requests) > 0 {
		if err := bot.PostMessages(ctx, requests, d.config.ClientID, d.config.SessionID, d.config.SessionKey); err != nil {
			log.WithError(err).Error("post messages")
			return err
		}
	}

	if err := d.commands.Deletes(ctx, cmds); err != nil {
		log.WithError(err).Error("delete commands")
		return err
	}

	return nil
}

func (d *Deliver) handleCommand(ctx context.Context, cmd *core.Command) ([]*bot.MessageRequest, error) {
	c, err := d.prepareContext(ctx, cmd)
	if err != nil {
		return nil, err
	}

	var requests []*bot.MessageRequest

	// 设置语言
	switch cmd.Action {
	case core.ActionSwitchChinese, core.ActionSwitchEnglish:
		if cmd.Action != c.user.Language {
			c.user.Language = cmd.Action
			if err := d.users.Update(ctx, c.user); err != nil {
				return nil, fmt.Errorf("update user failed: %w", err)
			}
		}

		req := c.languageSwitched(ctx)
		requests = append(requests, req)
		return requests, nil
	}

	// 还没有设置语言
	if c.user.Language == "" {
		req := c.selectLanguage(ctx, cmd)
		requests = append(requests, req)
		return requests, nil
	}

	switch cmd.Action {
	case core.ActionUsage:
		requests = append(requests, c.showUsage(ctx))
	case core.ActionRequestCoin:
		// 每个小时只能领取一次
		req := &core.Transfer{
			TraceID:    uuid.Modify(cmd.UserID, cmd.CreatedAt.Truncate(time.Hour).Format(time.RFC3339)),
			OpponentID: cmd.UserID,
			AssetID:    d.config.CoinAsset,
			Amount:     "10",
			Memo:       "from blockquiz",
		}
		if err := d.wallets.CreateTransfer(ctx, req); err != nil {
			return nil, fmt.Errorf("create transfer failed: %w", err)
		}
	case core.ActionRandomCourse:
		course, err := d.pickRandomCourse(ctx, c.user)
		if err != nil {
			if store.IsErrNotFound(err) {
				break
			}
			return nil, fmt.Errorf("pick random course failed: %w", err)
		}

		c.course = course
		requests = append(requests, c.showCourseContent(ctx))
		requests = append(requests, c.showCourseButtons(ctx))
	case core.ActionShowCourse:
		if c.course == nil {
			break
		}

		requests = append(requests, c.showCourseContent(ctx))
		requests = append(requests, c.showCourseButtons(ctx))
	case core.ActionShowQuestion:
		if c.question == nil {
			break
		}

		requests = append(requests, c.showQuestionContent(ctx))
		requests = append(requests, c.showQuestionChoiceButtons(ctx))
	case core.ActionAnswerQuestion:
		if c.question == nil {
			break
		}

		right := c.question.Answer == cmd.Answer
		requests = append(requests, c.answerFeedback(ctx, right))

		if _, ok := c.course.Question(cmd.Question + 1); ok {
			requests = append(requests, c.showNextQuestionButton(ctx, cmd.Question+1))
		} else {
			requests = append(requests, c.showFinishCourse(ctx))
			if nextCourse, err := d.courses.FindNext(ctx, c.course.ID); err == nil {
				requests = append(requests, c.showNextCourseButton(ctx, nextCourse))
			}
		}
	default:
		logger.FromContext(ctx).Warnf("unknown action %s", cmd.Action)
	}

	return requests, nil
}

func (d *Deliver) pickRandomCourse(ctx context.Context, user *core.User) (*core.Course, error) {
	list, err := d.courses.ListAll(ctx, user.Language)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, store.ErrNotFound
	}

	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})

	return list[0], nil
}
