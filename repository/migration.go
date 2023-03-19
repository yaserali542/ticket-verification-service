package repository

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/yaserali542/ticket-verification-service/models"
)

func DropVerifiedQRCodesTable(db *gorm.DB) error {
	return db.DropTable(&models.VerifiedQRCodes{}).Error
}
func MigrateVerifiedQRCodesTable(db *gorm.DB) error {
	if !db.HasTable(&models.VerifiedQRCodes{}) {
		sqlStr := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
		tx := db.Exec(sqlStr)

		if tx.Error != nil {
			fmt.Println(tx.Error)
		} else {
			fmt.Println("create extension ran successfully")
		}
	}

	if err := db.AutoMigrate(&models.VerifiedQRCodes{}); err.Error != nil {
		return err.Error
	}
	if !db.HasTable(&models.VerifiedQRCodes{}) {
		return errors.New("table doesn't exist")
	}
	return nil
}
