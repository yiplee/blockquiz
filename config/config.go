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
		DB      db.Config `json:"db,omitempty"`
		Bot     Bot       `json:"bot,omitempty"`
		Course  Course    `json:"course,omitempty"`
		I18n    I18n      `json:"i18n,omitempty"`
		Deliver Deliver   `json:"deliver,omitempty"`
	}

	Bot struct {
		ClientID   string `json:"client_id,omitempty"`
		SessionID  string `json:"session_id,omitempty"`
		SessionKey string `json:"session_key,omitempty"`
		PinToken   string `json:"pin_token,omitempty"`
		Pin        string `json:"pin,omitempty"`
	}

	Course struct {
		// 课程文件路径
		Path string `json:"path,omitempty"`

		// 答题用的币的 asset id
		CoinAsset  string          `json:"coin_asset,omitempty"`
		CoinAmount decimal.Decimal `json:"coin_amount,omitempty"`

		// 答对奖励
		RewardAssetID string          `json:"reward_asset_id,omitempty"`
		RewardAmount  decimal.Decimal `json:"reward_amount,omitempty"`
	}

	I18n struct {
		Path string `json:"path,omitempty"`
	}

	Deliver struct {
		ButtonColor   string `json:"button_color,omitempty"`
		BlockDuration int64  `json:"block_duration,omitempty"` // 秒
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

	defaultConfig(&cfg)
	return &cfg, nil
}

func defaultConfig(cfg *Config) {
	if cfg.Deliver.BlockDuration == 0 {
		cfg.Deliver.BlockDuration = 60*60 + 1
	}

	if cfg.Deliver.ButtonColor == "" {
		cfg.Deliver.ButtonColor = "#11A7F7"
	}
}
