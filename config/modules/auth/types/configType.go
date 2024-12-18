package types

import (
	"time"

	basictypes "github.com/BIQ-Cat/easyserver/config/types"
)

type Config struct {
	basictypes.BasicConfig
	Create              CreateConfig           `json:"create"`          // Configure data required to create account
	OTPLength           int                    `json:"otpLength"`       // How long one-time password should be
	Verify              VerificationConfig     `json:"verify"`          // Configure verification
	RestorePassword     RestorePasswordConfig  `json:"restorePassword"` // Restore password configuration
	OAuthConfigsEnabled map[string]OAuthConfig `json:"-"`               // All OAuth2 configurations enabled for the server. Here, the key should be a handler name passing in request.
}

func (c Config) GetJSONPath() string {
	return "./EasyConfig/modules/auth.json"
}

type CreateConfig struct {
	DisableUsername bool             `json:"disableUsername"` // Whether disable username. Should be false if neither email nor phone required
	Phone           EmailPhoneConfig `json:"phone"`           // Configure phone number usage
	Email           EmailPhoneConfig `json:"email"`           // Configure email address usage
	RequireEither   bool             `json:"requireEither"`   // Works when both Email.Require and Phone.Require are true. If true, either phone or email is required. If false, both of them
}

func (cfg CreateConfig) IsEmailRequired(phone string) bool {
	return cfg.Email.isRequired(!cfg.RequireEither, phone, cfg.Phone)
}
func (cfg CreateConfig) IsPhoneRequired(email string) bool {
	return cfg.Phone.isRequired(!cfg.RequireEither, email, cfg.Email)
}

type VerificationConfig struct {
	Require          bool          `json:"require"`          // Disables account before verification
	SetPasswordAfter bool          `json:"setPasswordAfter"` // If enabled, password cannot be set before verification. So, accont is saved into database only after verification and setting password (it is required before using account)
	EmailSubject     string        `json:"emailSubject"`     // Subject for email sending
	ResendTimer      time.Duration `json:"resendTimer"`      // How many time should pass before re-sending?
	TokenLifetime    time.Duration `json:"tokenLifetime"`    // How long token will be available?
}

type EmailPhoneConfig struct {
	Require    bool `json:"require"`    // Make email / phone required.
	UseAsLogin bool `json:"useAsLogin"` // Enables logging in by email / phone
}

type RestorePasswordConfig struct {
	EmailSubject  string        `json:"emailSubject"`  // Subject for email sending
	ResendTimer   time.Duration `json:"resendTimer"`   // How many time should pass before re-sending?
	TokenLifetime time.Duration `json:"tokenLifetime"` // How long token will be available?
}

func (cfg EmailPhoneConfig) isRequired(requireBoth bool, other string, otherCfg EmailPhoneConfig) bool {
	if !cfg.Require { // Not required
		return false
	}
	if requireBoth {
		return true
	}
	// Either phone or email? One only?
	return otherCfg.Require || other == ""
}
