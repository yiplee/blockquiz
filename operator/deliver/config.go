package deliver

type Config struct {
	ClientID   string `valid:"uuid,required"`
	SessionID  string `valid:"uuid,required"`
	SessionKey string `valid:"required"`

	ButtonColor string `valid:"required"`
	// 用户自己创建的 task 回答错误的话需要等待的时长，单位秒
	BlockDuration int64
	QuestionCount int `valid:"required"`
}
