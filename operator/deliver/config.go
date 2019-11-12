package deliver

import (
	"time"
)

type Config struct {
	ClientID   string `valid:"uuid,required"`
	SessionID  string `valid:"uuid,required"`
	SessionKey string `valid:"required"`

	ButtonColor   string `valid:"required"`
	BlockDuration time.Duration
}
