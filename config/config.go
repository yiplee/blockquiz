package config

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"github.com/yiplee/blockquiz/db"
)

type (
	Config struct {
		Bot    Bot       `json:"bot,omitempty"`
		Lesson Lesson    `json:"lesson,omitempty"`
		I18n   I18n      `json:"i18n,omitempty"`
		DB     db.Config `json:"db,omitempty"`
	}

	Bot struct {
		ClientID   string `json:"client_id,omitempty"`
		SessionID  string `json:"session_id,omitempty"`
		SessionKey string `json:"session_key,omitempty"`
		PinToken   string `json:"pin_token,omitempty"`
		Pin        string `json:"pin,omitempty"`
	}

	Lesson struct {
		Path string `json:"path,omitempty"`
		// 在发送问题的时候是否打乱答案的顺序，同样的问题不同用户看到的选项顺序不一样
		ShuffleChoices bool `json:"shuffle_choices,omitempty"`
		// 答题用的币的 asset id
		Coin   string          `json:"coin,omitempty"`
		Amount decimal.Decimal `json:"amount,omitempty"`
	}

	I18n struct {
		Path string `json:"path,omitempty"`
	}
)

func Load(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("open config file at %s failed: %w", configPath, err)
	}
	defer f.Close()

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(f); err != nil {
		return nil, fmt.Errorf("viper read config failed: %w", err)
	}

	data, err := jsoniter.Marshal(viper.AllSettings())
	if err != nil {
		return nil, fmt.Errorf("marshal config data failed: %w", err)
	}

	var cfg Config
	if err := jsoniter.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config data failed: %w", err)
	}

	return &cfg, nil
}
