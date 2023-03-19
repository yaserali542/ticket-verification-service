package repository

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/yaserali542/ticket-verification-service/models"
)

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) GetVerificationQRData(signature string) (*models.VerifiedQRCodes, bool, error) {

	var data models.VerifiedQRCodes
	db := r.Db.First(&data, "digital_signature =?", signature)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, true, nil
		} else {
			return nil, true, db.Error
		}
	}
	return &data, false, nil

}
func (r *Repository) CreateVerificationQRData(booking *models.VerifiedQRCodes) error {
	return r.Db.Create(booking).Error
}
