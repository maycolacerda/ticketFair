package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BeforeCreateInterface interface {
	BeforeCreate(tx *gorm.DB) (err error)
}

func (User *User) BeforeCreate(tx *gorm.DB) (err error) {

	User.UserID = uuid.New().String()

	User.Password, err = hashPassword(User.Password)
	if err != nil {
		return err
	}
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
	MerchantRep.Password, err = hashPassword(MerchantRep.Password)
	if err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
