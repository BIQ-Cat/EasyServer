package models

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/BIQ-Cat/easyserver/db"

	"github.com/BIQ-Cat/easyserver/config"
	moduleConfig "github.com/BIQ-Cat/easyserver/modules/auth/config"
	"github.com/BIQ-Cat/easyserver/utils"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

func (a *Account) SendEmailOTP(email string, isVerification bool, host string) (map[string]interface{}, error) {
	type emailData struct {
		URL       string
		FirstName string
		Subject   string
		Token     string
	}

	if a.Verified {
		return utils.Message(false, "Account is already verified"), nil
	}

	if !moduleConfig.Config.Create.Email.Require {
		return utils.Message(false, "Email verification is disabled on this server"), nil
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
		URL:       host + "/auth/" + controller + "?token=" + base64.StdEncoding.EncodeToString([]byte(otp)) + "&visual=1",
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

func VerifyAccount(otp []byte) (map[string]interface{}, error) {
	acc, ok, err := findUserByField("verification_otp",
		base64.StdEncoding.EncodeToString(
			pbkdf2.Key(otp, []byte(config.EnvConfig.OTPPassword), moduleConfig.PBKDF2Iter, moduleConfig.PBKDF2Length, sha256.New),
		),
	)
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
	acc.VerificationOTP = ""
	err = db.GetDB().Save(acc).Error
	if err != nil {
		return nil, err
	}

	acc.Password = ""

	err = acc.generateToken()
	if err != nil {
		return nil, err
	}

	resp := utils.Message(true, "Account has been verified")
	resp["account"] = acc
	return resp, nil
}

func ResetPassword(otp, password []byte) (map[string]interface{}, error) {
	acc, ok, err := findUserByField("restore_password_otp",
		base64.StdEncoding.EncodeToString(
			pbkdf2.Key(otp, []byte(config.EnvConfig.OTPPassword), moduleConfig.PBKDF2Iter, moduleConfig.PBKDF2Length, sha256.New),
		),
	)
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

	acc.RestorePasswordOTP = ""
	resp, err := acc.ChangePassword(password) // Saves account internally
	if err != nil || !utils.GetStatus(resp) {
		return resp, err
	}

	acc.Password = ""

	err = acc.generateToken()
	if err != nil {
		return nil, err
	}

	resp["account"] = acc
	return resp, nil
}

func (a *Account) ChangePassword(password []byte) (map[string]interface{}, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	a.Password = string(hashedPassword)
	err = db.GetDB().Save(a).Error
	if err != nil {
		return nil, err
	}

	return utils.Message(true, "Changed password"), nil
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
		newOTP := base64.StdEncoding.EncodeToString(pbkdf2.Key(otpToken, salt, 4096, 32, sha256.New))

		_, ok, err := findUserByField(fieldName, newOTP)
		if err != nil {
			return "", nil, err
		}

		if ok {
			otpToken = utils.RandomText(moduleConfig.Config.OTPLength)
			continue
		}

		*otp = newOTP
		*otpSet = time.Now()
		break
	}

	err := db.GetDB().Save(a).Error
	if err != nil {
		return "", nil, err
	}

	return string(otpToken), nil, nil
}
