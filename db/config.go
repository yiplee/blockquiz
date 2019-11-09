package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Config struct {
	Dialect  string `json:"dialect,omitempty"` // mysql,postgres,sqlite
	Host     string `json:"host,omitempty"`    // if Dialect is `sqlite`, host should be db file path
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Debug    bool   `json:"debug,omitempty"`
}

func SqliteInMemory() Config {
	return Config{
		Dialect: "sqlite",
		Host:    ":memory:",
	}
}

func Open(cfg Config) (*DB, error) {
	var uri string
	switch cfg.Dialect {
	case "mysql":
		uri = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=True&charset=utf8mb4",
			cfg.User,
			cfg.Password,
			"tcp",
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
	case "postgres":
		uri = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Database,
			cfg.Password,
		)
	case "sqlite":
		uri = cfg.Host
	default:
		return nil, fmt.Errorf("unkonow db dialect: %s", cfg.Dialect)
	}

	db, err := gorm.Open(cfg.Dialect, uri)
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		db = db.Debug()
	}

	db.DB().SetMaxIdleConns(10)
	return &DB{
		write: db,
		read:  db,
	}, nil
}

func MustOpen(cfg Config) *DB {
	db, err := Open(cfg)
	if err != nil {
		panic(fmt.Errorf("open db failed: %w", err))
	}

	return db
}
