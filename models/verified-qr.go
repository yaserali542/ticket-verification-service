package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type VerifiedQRCodes struct {
	ID               uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"-"`
	BookingId        uuid.UUID `gorm:"type:uuid;not null;index" json:"booking_id"`
	DigitalSignature string    `gorm:"not null;index" json:"signature"`
	UserName         string    `json:"user_name"`
	VerifierUserName string    `json:"verfier_user_name"`
	CreatedAt        time.Time `json:"verified_at"`
	UpdatedAt        time.Time
}
