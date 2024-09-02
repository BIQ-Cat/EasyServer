package models

import (
	"net/mail"
	"regexp"

	"github.com/jinzhu/gorm"

	// Internals
	"github.com/BIQ-Cat/easyserver/internal/db"
	"github.com/BIQ-Cat/easyserver/internal/utils"

	// Configuration
	moduleConfig "github.com/BIQ-Cat/easyserver/config/modules/auth"
)

func (a *Account) Validate() (msg map[string]interface{}, ok bool) {
	ok = true
	msg = nil

	if a.Email == "" && moduleConfig.Config.Create.IsEmailRequired(a.Phone) {
		return utils.Message(false, "Email address is required"), false
	} else if a.Email != "" {
		msg, ok = a.validateEmail()
	}

	if a.Phone == "" && moduleConfig.Config.Create.IsPhoneRequired(a.Email) {
		return utils.Message(false, "Phone number is required"), false
	} else if a.Phone != "" {
		msg, ok = a.validatePhoneNumber()
	}

	if a.Password == "" && !moduleConfig.Config.Verify.SetPasswordAfter {
		return utils.Message(false, "Password is required"), false
	} else if a.Password != "" && !moduleConfig.Config.Verify.SetPasswordAfter {
		msg, ok = a.validatePassword()
	}

	if !ok {
		return
	}

	if !moduleConfig.Config.Create.DisableUsername {
		msg, ok = a.validateUsername()
	}
	return
}

func (a *Account) validateEmail() (map[string]interface{}, bool) {
	if _, err := mail.ParseAddress(a.Email); err != nil {
		return utils.Message(false, "Email address is incorrect"), false
	}

	temp := &Account{}
	err := db.GetDB().Table("accounts").Where("email = ?", a.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email == a.Email {
		return utils.Message(false, "Email address already in use by another user."), false
	}

	return nil, true
}

func (a *Account) validatePhoneNumber() (map[string]interface{}, bool) {
	if expr := regexp.MustCompile(`^\+[1-9]\d{1,14}$`); !expr.MatchString(a.Phone) {
		return utils.Message(false, "Phone number is incorrect"), false
	}

	temp := &Account{}
	err := db.GetDB().Table("accounts").Where("phone = ?", a.Phone).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.Message(false, "Connection error. Please retry"), false
	}
	if temp.Phone == a.Phone {
		return utils.Message(false, "Phone number already in use by another user."), false
	}
	return nil, true
}

func (a *Account) validatePassword() (map[string]interface{}, bool) {
	if len(a.Password) < 6 {
		return utils.Message(false, "Password is too weak"), false
	}
	return nil, true
}

func (a *Account) validateUsername() (map[string]interface{}, bool) {
	if a.Username == "" {
		return utils.Message(false, "Username is required"), false
	}

	temp := &Account{}
	err := db.GetDB().Table("accounts").Where("username = ?", a.Username).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.Message(false, "Connection error. Please retry"), false
	}
	if temp.Username == a.Username {
		return utils.Message(false, "Username already in use by another user."), false
	}

	return nil, true
}
