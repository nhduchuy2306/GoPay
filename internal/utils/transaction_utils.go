package utils

import "gorm.io/gorm"

type TxMode string

const (
	Required    TxMode = "REQUIRED"
	RequiresNew TxMode = "REQUIRES_NEW"
)

func WithTx(db *gorm.DB, mode TxMode, fn func(tx *gorm.DB) error) error {
	switch mode {
	case Required:
		return txRequired(db, fn)
	case RequiresNew:
		return txRequiresNew(db, fn)
	default:
		return txRequired(db, fn)

	}
}

func txRequired(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	if isTx(db) {
		return fn(db)
	}
	tx := db.Begin()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func txRequiresNew(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func isTx(db *gorm.DB) bool {
	return db.Statement != nil && db.Statement.ConnPool != nil
}
