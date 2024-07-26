package settings

const (
	CREATE_USERNAME_ONLY_REQUIRED = iota - 1
	CREATE_EMAIL_ONLY_REQUIRED
	CREATE_PHONE_ONLY_REQUIRED
	CREATE_BOTH_REQUIRED
	CREATE_EITHER_REQUIRED
)

type createSettings struct {
	RequireData                   int
	RequireVerification           bool
	SetPasswordBeforeVerification bool
	HasUsername                   bool
}

var Create = createSettings{
	CREATE_EITHER_REQUIRED,
	true,
	true,
	true,
}

func (c createSettings) IsEmailRequired(phoneNumber string) bool {
	return c.RequireData == CREATE_BOTH_REQUIRED || c.RequireData == CREATE_EMAIL_ONLY_REQUIRED || phoneNumber == "" && c.RequireData == CREATE_EITHER_REQUIRED
}

func (c createSettings) IsPhoneRequired(emailAdress string) bool {
	return c.RequireData == CREATE_BOTH_REQUIRED || c.RequireData == CREATE_PHONE_ONLY_REQUIRED || emailAdress == "" && c.RequireData == CREATE_EITHER_REQUIRED
}

func (c createSettings) Validate() bool {
	return c.RequireData == CREATE_USERNAME_ONLY_REQUIRED && !c.HasUsername || !c.RequireVerification && !c.SetPasswordBeforeVerification
}
