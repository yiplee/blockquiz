package hub

type Config struct {
	ClientID   string `valid:"uuid,required"`
	SessionID  string `valid:"uuid,required"`
	SessionKey string `valid:"required"`
}
