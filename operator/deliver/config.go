package deliver

import (
	"github.com/shopspring/decimal"
)

type Config struct {
	ClientID   string `valid:"uuid,required"`
	SessionID  string `valid:"uuid,required"`
	SessionKey string `valid:"required"`

	CoinAsset  string          `valid:"uuid,required"`
	CoinAmount decimal.Decimal `valid:"required"`

	ButtonColor string `valid:"required"`
}
