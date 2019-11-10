package deliver

import (
	"context"
	"encoding/base64"
	"net/url"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/fox-one/pkg/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/yiplee/blockquiz/core"
)

func (d *Deliver) buttonAction(ctx context.Context, cmds []*core.Command) string {
	uri, _ := url.Parse("mixin://pay")
	query := uri.Query()
	query.Set("opponent_id", d.config.ClientID)
	query.Set("asset_id", d.config.CoinAsset)
	query.Set("amount", d.config.CoinAmount.Truncate(8).String())
	query.Set("trace_id", uuid.New())
	query.Set("memo", d.parser.Encode(ctx, cmds))
	return uri.String()
}

type button struct {
	Label  string `json:"label,omitempty"`
	Color  string `json:"color,omitempty"`
	Action string `json:"action,omitempty"`
}

func (d *Deliver) selectLanguage(ctx context.Context, user *core.User) *bot.MessageRequest {
	req := &bot.MessageRequest{
		ConversationId: bot.UniqueConversationId(user.MixinID, d.config.ClientID),
		RecipientId:    user.MixinID,
		Category:       "APP_BUTTON_GROUP",
	}

	var buttons []button
	for _, lang := range []string{core.ActionSwitchEnglish, core.ActionSwitchChinese} {
		l := localizer.WithLanguage(d.localizer, lang)
		cmds := []*core.Command{{Action: lang}}

		buttons = append(buttons, button{
			Label:  l.MustLocalize("select_language"),
			Color:  d.config.ButtonColor,
			Action: d.parser.Encode(ctx, cmds),
		})
	}

	data, _ := jsoniter.Marshal(buttons)
	req.Data = base64.StdEncoding.EncodeToString(data)
	return req
}
