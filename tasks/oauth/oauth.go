package oauth

import (
	"time"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
	log "github.com/go-pkgz/lgr"
	"golang.org/x/oauth2"
)

func Init() *auth.Service {
	options := auth.Opts{
		SecretReader: token.SecretFunc(func(_ string) (string, error) { // secret key for JWT, ignores aud
			return "secret", nil
		}),
		TokenDuration:     time.Minute,                                 // short token, refreshed automatically
		CookieDuration:    time.Hour * 24,                              // cookie fine to keep for long time
		DisableXSRF:       true,                                        // don't disable XSRF in real-life applications!
		Issuer:            "aTES",                                      // part of token, just informational
		URL:               "http://localhost:8082",                     // base url of the protected service
		AvatarStore:       avatar.NewLocalFS("/tmp/demo-auth-service"), // stores avatars locally
		AvatarResizeLimit: 200,                                         // resizes avatars to 200x200
		Logger:            log.Default(),                               // optional logger for auth library
	}
	service := auth.NewService(options)
	service.AddCustomProvider("custom123", auth.Client{Cid: "cid", Csecret: "csecret"}, provider.CustomHandlerOpt{
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:9097/authorize",
			TokenURL: "http://localhost:9097/access_token",
		},
		InfoURL: "http://localhost:9097/user",
		MapUserFn: func(data provider.UserData, _ []byte) token.User {
			userInfo := token.User{
				ID:      data.Value("id"),
				Name:    data.Value("name"),
				Picture: data.Value("picture"),
			}
			return userInfo
		},
	})

	return service
}
