package models

import (
	"fmt"
	"log"

	"github.com/BIQ-Cat/easyserver/config"
	"github.com/BIQ-Cat/easyserver/db"
	moduleConfig "github.com/BIQ-Cat/easyserver/modules/auth/config"
	"github.com/BIQ-Cat/easyserver/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func (a *Account) Create() (map[string]interface{}, error) {
	if msg, ok := a.Validate(); !ok {
		return msg, nil
	}

	a.Verifyed = false

	if moduleConfig.Config.Create.SetPasswordBeforeVerification {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		a.Password = string(hashedPassword)

		db.GetDB().Create(a)

		if a.ID <= 0 {
			return utils.Message(false, "Failed to create account, connection error."), nil
		}
	}

	err := a.generateToken()
	if err != nil {
		return nil, err
	}

	a.Password = ""

	response := utils.Message(true, "Account has been created")
	response["account"] = a
	return response, nil
}

func Login(username, password string) (map[string]interface{}, error) {
	account := &Account{}

	err := db.GetDB().Table("accounts").Where("username = ?", username).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(false, "Email address not found"), nil
		}
		if config.Config.Debug {
			return nil, err
		} else {
			log.Println(fmt.Errorf("ERROR: Login: Check username: %w", err))
			return utils.Message(false, "Connection error. Please retry"), nil
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return utils.Message(false, "Invalid login credentials. Please try again"), nil
		}
		return nil, err
	}

	account.Password = ""

	err = account.generateToken()
	if err != nil {
		return nil, err
	}

	resp := utils.Message(true, "Logged In")
	resp["account"] = account
	return resp, nil
}
