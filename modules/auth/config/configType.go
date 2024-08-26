package config

import "time"

type config struct {
	Create          createConfig          `json:"create"`          // Configure data required to create account
	OTPLength       int                   `json:"otpLength"`       // How long one-time password should be
	Verify          verificationConfig    `json:"verify"`          // Configure verification
	RestorePassword restorePasswordConfig `json:"restorePassword"` // Restore password configuration
	RewriteWithJSON bool                  `json:"-"`               // Enables JSON configuration. If it exists (path ./json/modules/auth.json), this configuration will be shadowed by JSON one
}

type createConfig struct {
	DisableUsername bool             `json:"disableUsername"` // Whether disable username. Should be false if neither email nor phone required
	Phone           emailPhoneConfig `json:"phone"`           // Configure phone number usage
	Email           emailPhoneConfig `json:"email"`           // Configure email address usage
	RequireEither   bool             `json:"requireEither"`   // Works when both Email.Require and Phone.Require are true. If true, either phone or email is required. If false, both of them
}

func (cfg createConfig) IsEmailRequired(phone string) bool {
	return cfg.Email.isRequired(!cfg.RequireEither, phone, cfg.Phone)
}
func (cfg createConfig) IsPhoneRequired(email string) bool {
	return cfg.Phone.isRequired(!cfg.RequireEither, email, cfg.Email)
}

type verificationConfig struct {
	Require          bool          `json:"require"`          // Disables account before verification
	SetPasswordAfter bool          `json:"setPasswordAfter"` // If enabled, password cannot be set before verification. So, accont is saved into database only after verification and setting password (it is required before using account)
	EmailSubject     string        `json:"emailSubject"`     // Subject for email sending
	ResendTimer      time.Duration `json:"resendTimer"`      // How many time should pass before re-sending?
	TokenLifetime    time.Duration `json:"tokenLifetime"`    // How long token will be available?
}

type emailPhoneConfig struct {
	Require    bool `json:"require"`    // Make email / phone required.
	UseAsLogin bool `json:"useAsLogin"` // Enables logging in by email / phone
}

type restorePasswordConfig struct {
	EmailSubject string        `json:"emailSubject"` // Subject for email sending
	ResendTimer  time.Duration `json:"resendTimer"`  // How many time should pass before re-sending?
}

func (cfg emailPhoneConfig) isRequired(requireBoth bool, other string, otherCfg emailPhoneConfig) bool {
	if !cfg.Require { // Not required
		return false
	}
	if requireBoth {
		return true
	}
	// Either phone or email? One only?
	return otherCfg.Require || other == ""
}
