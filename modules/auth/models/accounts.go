package models

import (
	"time"

	"github.com/BIQ-Cat/easyserver/config"
	"github.com/BIQ-Cat/easyserver/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
)

type Token struct {
	UserId   uint
	Verified bool
	jwt.RegisteredClaims
}

type Account struct {
	gorm.Model
	Email                    string    `json:"email,omitempty"`
	Password                 string    `json:"password"`
	Token                    string    `json:"token" gorm:"-:all"`
	Phone                    string    `json:"phone,omitempty"`
	Username                 string    `json:"username"`
	Verified                 bool      `json:"verified"`
	VerificationOTP          string    `json:"-"`
	RestorePasswordOTP       string    `json:"-"`
	TimeVerificationOTPSet   time.Time `json:"-"`
	TimeForgotPasswordOTPSet time.Time `json:"-"`
}

func GetUser(u uint) *Account {

	acc := &Account{}
	db.GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" {
		return nil
	}

	acc.Password = ""
	return acc
}

func (a *Account) generateToken() error {
	tk := &Token{UserId: a.ID, Verified: a.Verified}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(config.EnvConfig.TokenPassword))
	if err != nil {
		return err
	}
	a.Token = tokenString

	return nil
}
