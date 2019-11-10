package wallet

import (
	"context"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/go-number"
	"github.com/asaskevich/govalidator"
	"github.com/yiplee/blockquiz/core"
)

type Config struct {
	ClientID   string `valid:"uuid,required"`
	SessionID  string `valid:"uuid,required"`
	PinToken   string `valid:"required"`
	Pin        string `valid:"required"`
	SessionKey string `valid:"required"`
}

type walletSrv struct {
	cfg Config
}

func New(cfg Config) core.WalletService {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	return &walletSrv{cfg: cfg}
}

func (w *walletSrv) Transfer(ctx context.Context, req *core.Transfer) error {
	input := &bot.TransferInput{
		AssetId:     req.AssetID,
		RecipientId: req.OpponentID,
		Amount:      number.FromString(req.Amount),
		TraceId:     req.TraceID,
		Memo:        req.Memo,
	}

	return bot.CreateTransfer(ctx, input, w.cfg.ClientID, w.cfg.SessionID, w.cfg.SessionKey, w.cfg.Pin, w.cfg.PinToken)
}
