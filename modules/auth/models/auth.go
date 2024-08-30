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

	a.Verified = !moduleConfig.Config.Create.Email.Require && !moduleConfig.Config.Create.Phone.Require

	if !moduleConfig.Config.Verify.SetPasswordAfter {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		a.Password = string(hashedPassword)
	}
	db.GetDB().Create(a)

	if a.ID <= 0 {
		return utils.Message(false, "Failed to create account, connection error."), nil
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

func Login(login, password string) (map[string]interface{}, error) {
	account := &Account{}
	var err error
	fields := make([]string, 3)

	if !moduleConfig.Config.Create.DisableUsername {
		fields = append(fields, "username")
	}
	if moduleConfig.Config.Create.Email.UseAsLogin {
		fields = append(fields, "email")
	}
	if moduleConfig.Config.Create.Phone.UseAsLogin {
		fields = append(fields, "phone")
	}

	for _, field := range fields {
		if field == "" {
			continue
		}
		err = db.GetDB().Table("accounts").Where(field+" = ?", login).First(account).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			if config.Config.Debug {
				return nil, err
			} else {
				log.Println(fmt.Errorf("ERROR: Login: Check username: %w", err))
				return utils.Message(false, "Connection error. Please retry"), nil
			}
		}
		break
	}

	if err != nil {
		return utils.Message(false, "Login not found"), nil
	}

	if !moduleConfig.Config.Verify.SetPasswordAfter || account.Verified {

		err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return utils.Message(false, "Invalid login credentials. Please try again"), nil
			}
			return nil, err
		}
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
