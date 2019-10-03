package oauth

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-session/session"

	"gopkg.in/oauth2.v3/server"

	"github.com/go-chi/chi"
)

type Router struct {
	ctx            context.Context
	Server         *server.Server
	SessionManager *session.Manager
}

func (router *Router) Run(chiRouter chi.Router) {
	chiRouter.Get("/authorize", router.Authorize)
	chiRouter.Post("/authorize", router.Authorize)

	chiRouter.Get("/token", router.Token)
	chiRouter.Post("/token", router.Token)
}

func (router *Router) Authorize(w http.ResponseWriter, r *http.Request) {
	store, err := router.SessionManager.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params, ok := store.Get("ReturnUri")
	if !ok {
		params = ""
	}

	form, err := url.ParseQuery(params.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(form) != 0 {
		r.Form = form
	}

	store.Delete("ReturnUri")
	err = store.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = router.Server.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (router *Router) Token(w http.ResponseWriter, r *http.Request) {
	err := router.Server.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func userAuthorization(router *Router) func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	return func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		sessionStore, err := router.SessionManager.Start(r.Context(), w, r)
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
