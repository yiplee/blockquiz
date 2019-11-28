package property

import (
	"context"
	"time"

	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
	"github.com/yiplee/blockquiz/property"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		if err := db.Update().AutoMigrate(Property{}).Error; err != nil {
			return err
		}

		return nil
	})
}

type Property struct {
	Key       string         `gorm:"size:64;PRIMARY_KEY"`
	Value     property.Value `gorm:"type:varchar(256)"`
	UpdatedAt time.Time      `gorm:"precision:6"`
}

type propertyStore struct {
	db *db.DB
}

func New(db *db.DB) core.PropertyStore {
	return &propertyStore{db: db}
}

func (s *propertyStore) Get(ctx context.Context, key string) (property.Value, error) {
	var p Property
	err := s.db.View().Where("`key` = ?", key).First(&p).Error
	if db.IsErrorNotFound(err) {
		err = nil
	}

	return p.Value, err
}

func (s *propertyStore) Save(ctx context.Context, key string, value interface{}) error {
	p := Property{
		Key:   key,
		Value: property.Parse(value),
	}

	tx := s.db.Update().Save(&p)
	if err := tx.Error; err != nil {
		return err
	}

	if tx.RowsAffected == 0 {
		return s.db.Update().Create(&p).Error
	}

	return nil
}

func (s *propertyStore) List(ctx context.Context) (map[string]property.Value, error) {
	var properties []Property
	if err := s.db.View().Find(&properties).Error; err != nil {
		return nil, err
	}

	values := make(map[string]property.Value, len(properties))
	for _, p := range properties {
		values[p.Key] = p.Value
	}

	return values, nil
}
