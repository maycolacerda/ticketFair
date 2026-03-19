package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BeforeCreateInterface interface {
	BeforeCreate(tx *gorm.DB) (err error)
}

func (User *User) BeforeCreate(tx *gorm.DB) (err error) {

	User.UserID = uuid.New().String()
	return nil
}

func (Profile *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	Profile.ProfileID = uuid.New().String()
	return nil
}

func (Merchant *Merchant) BeforeCreate(tx *gorm.DB) (err error) {
	Merchant.MerchantID = uuid.New().String()
	return nil
}

func (MerchantRep *MerchantRep) BeforeCreate(tx *gorm.DB) (err error) {
	MerchantRep.MerchantRepID = uuid.New().String()

	return nil
}
