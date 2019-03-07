package oauth

import (
	"context"
	"log"
	"net/http"

	"github.com/go-session/session"

	"gopkg.in/oauth2.v3"

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"

	"time"
)

type Oauth2Service struct {
	SessionManager *session.Manager
	TokenStore     oauth2.TokenStore
	ClientStore    *store.ClientStore
}

type OauthServer interface {
	HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) (err error)
	HandleTokenRequest(w http.ResponseWriter, r *http.Request) (err error)
	ValidationBearerToken(r *http.Request) (ti oauth2.TokenInfo, err error)
}

type oauth2Server struct {
	*server.Server
}

func NewOauthServer(oauth2Service *Oauth2Service) OauthServer {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(
		&manage.Config{
			AccessTokenExp:    time.Hour * 24 * 7,
			RefreshTokenExp:   time.Hour * 24 * 14,
			IsGenerateRefresh: true,
		},
	)

	manager.MapTokenStorage(oauth2Service.TokenStore)
	manager.MapClientStorage(oauth2Service.ClientStore)

	server := server.NewDefaultServer(manager)
	server.UserAuthorizationHandler = userAuthorization(oauth2Service)
	server.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Printf("[ERROR] Internal Error: %s", err.Error())
		return
	})
	server.SetResponseErrorHandler(func(re *errors.Response) {
		log.Printf("[ERROR] Response Error: %s", re.Error.Error())
	})

	return NewOauthServerWithServer(server)
}

func NewOauthServerWithServer(srv *server.Server) OauthServer {
	return &oauth2Server{Server: srv}
}

func userAuthorization(service *Oauth2Service) func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	return func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		log.Printf("[INFO] userAuthorization %s", r.URL)
		sessionStore, err := service.SessionManager.Start(context.Background(), w, r)
		if err != nil {
			return
		}

		uid, ok := sessionStore.Get("LoggedInUserID")
		if !ok {
			if r.Form == nil {
				r.ParseForm()
			}

			sessionStore.Set("ReturnUri", r.Form.Encode())
			sessionStore.Save()

			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusFound)
			return
		}
		userID = uid.(string)

		// Authorization for receiving a token
		sessionStore.Delete("LoggedInUserID")
		sessionStore.Save()

		return
	}
}
