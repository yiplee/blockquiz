package db

import (
	"github.com/hashicorp/go-multierror"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type DB struct {
	write *gorm.DB
	read  *gorm.DB
}

func (db *DB) Update() *gorm.DB {
	return db.write
}

func (db *DB) View() *gorm.DB {
	return db.read
}

func (db *DB) Debug() *DB {
	return &DB{
		write: db.write.Debug(),
		read:  db.read.Debug(),
	}
}

func (db *DB) Begin() *DB {
	tx := db.write.Begin()

	return &DB{
		write: tx,
		read:  db.read,
	}
}

func (db *DB) Rollback() error {
	return db.write.Rollback().Error
}

func (db *DB) Commit() error {
	return db.write.Commit().Error
}

func (db *DB) RollbackUnlessCommitted() {
	if err := db.write.RollbackUnlessCommitted().Error; err != nil {
		logrus.WithError(err).Error("DB: RollbackUnlessCommitted")
	}
}

func (db *DB) Tx(fn func(tx *DB) error) error {
	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) Ping() error {
	return db.write.DB().Ping()
}

func (db *DB) Close() error {
	var merr *multierror.Error

	if err := db.write.Close(); err != nil {
		merr = multierror.Append(merr, err)
	}

	if err := db.read.Close(); err != nil {
		merr = multierror.Append(merr, err)
	}

	return merr.ErrorOrNil()
}
