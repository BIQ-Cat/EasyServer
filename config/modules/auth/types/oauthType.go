package types

import "golang.org/x/oauth2"

type OAuthConfig struct {
	oauth2.Config
	UserInfo func(token string) ([]byte, error)
}
