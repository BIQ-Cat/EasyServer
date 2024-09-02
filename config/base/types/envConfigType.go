package types

type EnvConfig struct {
	DBName string `name:"DB_NAME"`
	DBPass string `name:"DB_PASS"`
	DBUser string `name:"DB_USER"`
	DBType string `name:"DB_TYPE"`
	DBHost string `name:"DB_HOST"`
	DBPort int    `name:"DB_PORT" default:"5432"`

	ServerPort int `name:"PORT" default:"8080"`

	TokenPassword string `name:"TOKEN_PASSWORD"`
	OTPPassword   string `name:"OTP_PASSWORD"`

	EmailFrom string `name:"EMAIL_FROM"`
	SMTPHost  string `name:"SMTP_HOST"`
	SMTPPass  string `name:"SMTP_PASS"`
	SMTPPort  int    `name:"SMTP_PORT"`
	SMTPUser  string `name:"SMTP_USER"`
}
