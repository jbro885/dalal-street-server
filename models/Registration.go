package models

import (
	"errors"

	"github.com/sirupsen/logrus"
)

var (
	VerificationKeyNotFoundError = errors.New("Verification key not found")
	TransactionError             = errors.New("Transaction failed")
	AlreadyVerifiedError         = errors.New("Account already verified")
)

type Registration struct {
	Id              uint32 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Email           string `gorm:"unique;not null" json:"email"`
	Name            string `gorm:"column:fullName;not null" json:"full_name"`
	UserName        string `gorm:"column:userName;not null" json:"user_name"`
	UserId          uint32 `gorm:"column:userId;not null" json:"user_id"`
	IsPragyan       bool   `gorm:"column:isPragyan;not null" json:"is_pragyan"`
	Password        string `gorm:"column:password;not null" json:"password"`
	Country         string `gorm:"not null" json:"country"`
	IsVerified      bool   `gorm:"column:isVerified;not null" json:"is_verified"`
	VerificationKey string `gorm:"column:verificationKey" json:"verification_key"`
	ReferralCodeID  uint32 `gorm:"column:referralCode;default:NULL" json:"referral_code"`
}

func (Registration) TableName() string {
	return "Registrations"
}

// VerifyAccount is used to handle logic after clicking the email verification
// link sent after registration
func VerifyAccount(verificationKey string) error {
	var l = logger.WithFields(logrus.Fields{
		"method":                "VerifyAccount",
		"param_verificationKey": verificationKey,
	})
	db := getDB()

	var registeredUser = Registration{
		VerificationKey: verificationKey,
	}

	err := db.Where("verificationKey = ?", verificationKey).First(&registeredUser).Error
	if err != nil {
		l.Errorf("No verification key %s", verificationKey)
		return VerificationKeyNotFoundError
	} else if registeredUser.IsVerified == true {
		l.Infof("Already Verified")
		return AlreadyVerifiedError
	} else {
		registeredUser.IsVerified = true
		if err := db.Save(registeredUser).Error; err != nil {
			l.Errorf("Error while modifiying IsVerified for key %s", verificationKey)
			registeredUser.IsVerified = false
		}
		return nil
	}
}
