package oauth

import (
	"net/http"

	"github.com/go-session/session"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/store"
)

type Service struct {
	SessionManager *session.Manager
	TokenStore     oauth2.TokenStore
	ClientStore    *store.ClientStore
}

func userAuthorization(service *Service) func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	return func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		sessionStore, err := service.SessionManager.Start(r.Context(), w, r)
		if err != nil {
			return
		}

		uid, ok := sessionStore.Get("LoggedInUserID")
		if !ok {
			if r.Form == nil {
				r.ParseForm()
			}

			sessionStore.Set("ReturnUri", r.Form.Encode())
			err = sessionStore.Save()
			if err != nil {
				return
			}

			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusFound)
			return
		}
		userID = uid.(string)

		// Authorization for receiving a token
		sessionStore.Delete("LoggedInUserID")
		err = sessionStore.Save()
		if err != nil {
			return
		}

		return
	}
}
