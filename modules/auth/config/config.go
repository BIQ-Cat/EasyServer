package config

import "time"

var Config = config{
	Create: createConfig{
		Phone: emailPhoneConfig{
			Require: false,
		},
		Email: emailPhoneConfig{
			Require:    true,
			UseAsLogin: true,
		},
		RequireEither: true,
	},
	Verify: verificationConfig{
		Require:       true,
		EmailSubject:  "Your account verification code",
		ResendTimer:   2 * time.Minute,
		TokenLifetime: 2 * time.Hour,
	},
	RestorePassword: restorePasswordConfig{
		EmailSubject: "Your password reset token",
		ResendTimer:  1 * time.Minute,
	},
	RewriteWithJSON: true,
	OTPLength:       6,
}
