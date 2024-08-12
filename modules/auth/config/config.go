package config

var Config = config{
	Create: createConfig{
		RequireData:                   CREATE_EITHER_REQUIRED,
		RequireVerification:           true,
		HasUsername:                   true,
		SetPasswordBeforeVerification: true,
	},
}
