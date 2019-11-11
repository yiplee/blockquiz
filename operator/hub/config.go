package hub

type Config struct {
	ClientID   string `valid:"uuid,required"`
	SessionID  string `valid:"uuid,required"`
	SessionKey string `valid:"required"`

	// 是否处理收款的 memo 作为用户输入的 command
	TransferCommand bool
}
