package hub

import (
	"context"
	"encoding/base64"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/fox-one/pkg/logger"
)

func (h *Hub) OnMessage(ctx context.Context, msg bot.MessageView, userId string) error {
	log := logger.FromContext(ctx).WithField("category", msg.Category)
	data, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		return nil
	}
	log.Info(string(data))
	return nil
}
