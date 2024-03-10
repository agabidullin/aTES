package oauth

import (
	"context"
	"fmt"

	"github.com/agabidullin/aTES/auth/model"

	"net/http"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	goauth2 "github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-pkgz/auth/provider"
	log "github.com/go-pkgz/lgr"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) {
	srv := initGoauth2Srv(db)
	sopts := provider.CustomServerOpt{
		URL:           "http://localhost:9097",
		L:             log.Default(),
		WithLoginPage: true,
	}
	// create custom provider and prepare params for handler
	prov := provider.NewCustomServer(srv, sopts)

	// Start server
	go prov.Run(context.Background())
}

// initialize go-oauth2/oauth2 server
func initGoauth2Srv(db *gorm.DB) *goauth2.Server {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("custom", []byte("00000000"), jwt.SigningMethodHS512))

	// client memory store
	clientStore := store.NewClientStore()
	err := clientStore.Set("cid", &models.Client{
		ID:     "cid",
		Secret: "csecret",
		Domain: "http://localhost:8082",
	})
	if err != nil {
		log.Printf("failed to set up a client store for go-oauth2/oauth2 server, %s", err)
	}
	err = clientStore.Set("cid1", &models.Client{
		ID:     "cid1",
		Secret: "csecret1",
		Domain: "http://localhost:8083",
	})
	if err != nil {
		log.Printf("failed to set up a client store for go-oauth2/oauth2 server, %s", err)
	}
	manager.MapClientStorage(clientStore)

	srv := goauth2.NewServer(goauth2.NewConfig(), manager)

	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (string, error) {
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		var account model.Account
		query := db.Where("Login = ? AND Password >= ?", username, password).First(&account)

		if query.Error != nil {
			return "", fmt.Errorf(query.Error.Error())
		}

		return fmt.Sprintf("custom123_%v", account.ID), nil
	})

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Printf("Internal Error: %s", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Printf("Response Error: %s", re.Error.Error())
	})

	return srv
}
