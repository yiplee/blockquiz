package hub

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/fox-one/pkg/logger"
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
		input = transfer.Memo
	}

	input = strings.ToLower(input)
	commands := strings.FieldsFunc(input, func(r rune) bool {
		return r == ';'
	})

	for idx, cmd := range commands {
		if args := strings.Fields(cmd); len(args) > 0 {
			traceID := uuid.Modify(msg.MessageId, strconv.Itoa(idx))
			if err := h.handleCommand(ctx, msg.UserId, traceID, args); err != nil {
				return err
			}
		}
	}

	return nil
}
