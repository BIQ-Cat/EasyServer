package models

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/BIQ-Cat/easyserver/db"

	"github.com/BIQ-Cat/easyserver/config"
	moduleConfig "github.com/BIQ-Cat/easyserver/modules/auth/config"
	"github.com/BIQ-Cat/easyserver/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/pbkdf2"
)

func (a *Account) SendEmailOTP(email string, isVerification bool, host string) (map[string]interface{}, error) {
	type emailData struct {
		URL       string
		FirstName string
		Subject   string
		Token     string
	}

	if a.Email != email {
		return utils.Message(false, "Incorrect email"), nil
	}

	otp, msg, err := a.setOTP(isVerification)
	if msg != nil || err != nil {
		return msg, err
	}

	var controller, subject, template string

	if isVerification {
		controller = "verify-recieve"
		subject = moduleConfig.Config.Verify.EmailSubject
		template = "verify.html"
	} else {
		controller = "reset-password"
		subject = moduleConfig.Config.RestorePassword.EmailSubject
		template = "reset.html"
	}

	data := emailData{
		URL:       host + "/auth/" + controller + "?token=" + base64.StdEncoding.EncodeToString([]byte(otp)) + "&style=visual",
		FirstName: a.Username,
		Subject:   subject,
		Token:     otp,
	}

	err = utils.SendEmail(a.Email, subject, &data, template)
	if err != nil {
		return utils.Message(false, "Error while sending E-mail: "+err.Error()), nil
	}

	return utils.Message(true, "OTP is sent on Email"), nil
}

func VerifyAccount(otp string) (map[string]interface{}, error) {
	acc, ok, err := findUserByField("verification_otp", string(
		pbkdf2.Key([]byte(otp), []byte(config.EnvConfig.OTPPassword), moduleConfig.PBKDF2Iter, moduleConfig.PBKDF2Length, sha256.New),
	))
	if err != nil {
		return nil, err
	}

	if !ok {
		return utils.Message(false, "No user with such token"), nil
	}

	if acc.Verified {
		return utils.Message(false, "Account is already verified"), nil
	}

	if time.Since(acc.TimeVerificationOTPSet) > moduleConfig.Config.Verify.TokenLifetime {
		return utils.Message(false, "Token has expired"), nil
	}

	acc.Verified = true
	err = db.GetDB().Save(acc).Error
	if err != nil {
		return nil, err
	}

	return utils.Message(true, "Account has been verified"), nil
}

func (a *Account) setOTP(isVerification bool) (string, map[string]interface{}, error) {
	var otp *string
	var otpSet *time.Time
	var resendTimer time.Duration
	var fieldName string

	if isVerification {
		otp = &a.VerificationOTP
		otpSet = &a.TimeVerificationOTPSet
		resendTimer = moduleConfig.Config.Verify.ResendTimer
		fieldName = "verification_otp"
	} else {
		otp = &a.RestorePasswordOTP
		otpSet = &a.TimeForgotPasswordOTPSet
		resendTimer = moduleConfig.Config.RestorePassword.ResendTimer
		fieldName = "restore_password_otp"
	}

	if time.Since(*otpSet) < resendTimer {
		return "", utils.Message(false, "Code has been sent earlier. Wait before re-sending a code"), nil
	}

	*otp = ""

	otpToken := utils.RandomText(moduleConfig.Config.OTPLength)
	salt := []byte(config.EnvConfig.OTPPassword)
	for {
		newOTP := string(pbkdf2.Key(otpToken, salt, 4096, 32, sha256.New))
		_, ok, err := findUserByField(fieldName, newOTP)
		if err != nil {
			return "", nil, err
		}
		if ok {
			otpToken = utils.RandomText(moduleConfig.Config.OTPLength)
			continue
		}

		*otp = newOTP
		break
	}

	err := db.GetDB().Save(a).Error
	if err != nil {
		return "", nil, err
	}

	return string(otpToken), nil, nil
}

func findUserByField(field, value string) (acc Account, ok bool, err error) {
	err = db.GetDB().Table("accounts").Where(field+" = ?", value).First(&acc).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	} else if err == nil {
		ok = true
	}

	return
}
