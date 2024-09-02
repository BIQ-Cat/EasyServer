package moduleConfig

import (
	"time"

	"github.com/BIQ-Cat/easyserver/config/modules/auth/types"
)

var Config = types.Config{
	Create: types.CreateConfig{
		Phone: types.EmailPhoneConfig{
			Require: false,
		},
		Email: types.EmailPhoneConfig{
			Require:    true,
			UseAsLogin: true,
		},
		RequireEither: true,
	},
	Verify: types.VerificationConfig{
		Require:       true,
		EmailSubject:  "Your account verification code",
		ResendTimer:   2 * time.Minute,
		TokenLifetime: 2 * time.Hour,
	},
	RestorePassword: types.RestorePasswordConfig{
		EmailSubject: "Your password reset token",
		ResendTimer:  1 * time.Minute,
	},
	RewriteWithJSON: true,
	OTPLength:       6,
}
