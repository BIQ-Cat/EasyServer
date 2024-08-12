package models

import (
	"os"

	"github.com/BIQ-Cat/easyserver/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type Token struct {
	UserId   uint
	Verified bool
	jwt.StandardClaims
}

type Account struct {
	gorm.Model
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
	Token    string `json:"token" gorm:"-:all"`
	Phone    string `json:"phone,omitempty"`
	Username string `json:"username"`
	Verifyed bool   `json:"verifyed"`
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
	tk := &Token{UserId: a.ID, Verified: a.Verifyed}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("H256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	if err != nil {
		return err
	}
	a.Token = tokenString

	return nil
}
