package funcs

import (
	"fmt"
	"io"
	"log"
	"net/http"

	config "github.com/BIQ-Cat/easyserver/config/base"
	"github.com/BIQ-Cat/easyserver/config/base/funcs"
	"github.com/BIQ-Cat/easyserver/config/modules/auth/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func OAuthGoogleConfig() types.OAuthConfig {
	if config.EnvConfig.OAuthGoogleClientID == funcs.EnvNoData || config.EnvConfig.OAuthGoogleClientSecret == funcs.EnvNoData {
		log.Fatalln("Google OAuth fields are not set!")
	}

	return types.OAuthConfig{
		Config: oauth2.Config{
			ClientID:     config.EnvConfig.OAuthGoogleClientID,
			ClientSecret: config.EnvConfig.OAuthGoogleClientSecret,
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint: google.Endpoint,
		},
		UserInfo: func(token string) ([]byte, error) {
			resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
			if err != nil {
				return nil, fmt.Errorf("failed getting user info: %v", err)
			}

			content, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed getting user info: %v", err)
			}
			return content, nil
		},
	}
}
