package hub

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/mq"
	"github.com/fox-one/pkg/number"
	"github.com/fox-one/pkg/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/thirdparty/bot-api-go-client"
)

func (h *Hub) OnMessage(ctx context.Context, msg bot.MessageView, userId string) error {
	log := logger.FromContext(ctx).WithField("category", msg.Category)

	if msg.UserId == "00000000-0000-0000-0000-000000000000" {
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		log.WithError(err).Warn("decode blaze message data")
		return nil
	}

	var (
		input  string
		source string
	)

	switch {
	case msg.Category == "PLAIN_TEXT":
		input = string(data)
		source = core.CommandSourcePlainText
	case msg.Category == "SYSTEM_ACCOUNT_SNAPSHOT" && h.config.TransferCommand:
		var transfer bot.TransferView
		_ = jsoniter.Unmarshal(data, &transfer)
		if amount := number.Decimal(transfer.Amount); !amount.IsPositive() {
			return nil
		}

		input = strings.ReplaceAll(transfer.Memo, "+", " ")
		source = core.CommandSourceSnapshot
	default:
		input = msg.Category
		source = core.CommandSourcePlainText
	}

	log.WithField("src", source).Debugf("parse input: %s", input)
	cmds, err := h.parser.Parse(ctx, input)
	if err != nil {
		return nil
	}

	traceID := msg.MessageId
	for idx, cmd := range cmds {
		if idx > 0 {
			traceID = uuid.Modify(traceID, strconv.Itoa(idx))
		}

		cmd.CreatedAt = msg.CreatedAt
		cmd.TraceID = traceID
		cmd.UserID = msg.UserId
		cmd.Source = source

		if err := h.handleCmd(ctx, cmd); err != nil {
			log.WithError(err).Error("pub cmd")
			return err
		}
	}

	return nil
}

func (h *Hub) handleCmd(ctx context.Context, cmd *core.Command) error {
	data, _ := jsoniter.MarshalToString(cmd)
	return h.pub.Publish(ctx, data, &mq.PublishOption{
		GroupID:   cmd.UserID,
		TraceID:   cmd.TraceID,
		ExpiredAt: time.Now().Add(time.Hour),
	})
}
