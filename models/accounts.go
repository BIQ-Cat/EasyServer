package models

import (
	"net/mail"
	"regexp"

	"github.com/BIQ-Cat/easyserver/routes/auth/settings"
	"github.com/BIQ-Cat/easyserver/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type Account struct {
	gorm.Model
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
	Token    string `json:"token" gorm:"-:all"`
	Phone    string `json:"phone,omitempty"`
	Username string `json:"username"`
}

func init() {
	modelsList = append(modelsList, Account{})
}

func (a *Account) Validate() (map[string]interface{}, bool) {
	if a.Email == "" && settings.Create.IsEmailRequired(a.Phone) {
		return utils.Message(false, "Email address is required"), false
	} else if a.Email != "" {
		if _, err := mail.ParseAddress(a.Email); err != nil {
			return utils.Message(false, "Email address is incorrect"), false
		}

		temp := &Account{}
		err := GetDB().Table("accounts").Where("email = ?", a.Email).First(temp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return utils.Message(false, "Connection error. Please retry"), false
		}
		if temp.Email == a.Email {
			return utils.Message(false, "Email address already in use by another user."), false
		}
	}

	if a.Phone == "" && settings.Create.IsEmailRequired(a.Email) {
		return utils.Message(false, "Phone number is required"), false
	} else if a.Phone != "" {
		if expr := regexp.MustCompile(`^\+[1-9]\d{1,14}$`); !expr.MatchString(a.Phone) {
			return utils.Message(false, "Phone number is incorrect"), false
		}

		temp := &Account{}
		err := GetDB().Table("accounts").Where("phone = ?", a.Phone).First(temp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return utils.Message(false, "Connection error. Please retry"), false
		}
		if temp.Phone == a.Phone {
			return utils.Message(false, "Phone number already in use by another user."), false
		}
	}

	if a.Password == "" && settings.Create.SetPasswordBeforeVerification {
		return utils.Message(false, "Password is required"), false
	} else if a.Password != "" {
		if len(a.Password) < 6 {
			return utils.Message(false, "Password is too weak"), false
		}
	}

	if a.Username == "" && settings.Create.HasUsername {
		return utils.Message(false, "Username is required"), false
	} else if a.Username != "" {
		temp := &Account{}
		err := GetDB().Table("accounts").Where("username = ?", a.Phone).First(temp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return utils.Message(false, "Connection error. Please retry"), false
		}
		if temp.Username == a.Username {
			return utils.Message(false, "Username already in use by another user."), false
		}
	} else if a.Email != "" {
		a.Username = a.Email
	} else if a.Phone != "" {
		a.Username = a.Email
	}
	return nil, true
}
