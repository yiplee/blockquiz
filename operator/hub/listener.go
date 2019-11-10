package hub

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/number"
	"github.com/fox-one/pkg/uuid"
	jsoniter "github.com/json-iterator/go"
)

func (h *Hub) OnMessage(ctx context.Context, msg bot.MessageView, userId string) error {
	log := logger.FromContext(ctx).WithField("category", msg.Category)

	data, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		log.WithError(err).Warn("decode blaze message data")
		return nil
	}

	var input string
	switch msg.Category {
	case "PLAIN_TEXT":
		input = string(data)
	case "SYSTEM_ACCOUNT_SNAPSHOT":
		var transfer bot.TransferView
		_ = jsoniter.Unmarshal(data, &transfer)
		if amount := number.Decimal(transfer.Amount); !amount.IsPositive() {
			return nil
		}

		input = strings.ReplaceAll(transfer.Memo, "+", " ")
	}

	cmds, err := h.parser.Parse(ctx, input)
	if err != nil {
		return nil
	}

	traceID := msg.MessageId
	for idx, cmd := range cmds {
		if idx > 0 {
			traceID = uuid.Modify(traceID, strconv.Itoa(idx))
		}

		cmd.TraceID = traceID
		cmd.UserID = msg.UserId
		if err := h.commands.Create(ctx, cmd); err != nil {
			log.WithError(err).Error("create command")
			return err
		}
	}

	return nil
}
