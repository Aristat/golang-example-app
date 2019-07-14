package oauth

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"

	"github.com/go-session/session"
)

var (
	invalidReturnUri = errors.New("10003 returnUri is not valid")
)

type Routers struct {
	SessionManager *session.Manager
	IServer
}

func (service *Routers) Run(router *chi.Mux) {
	router.Get("/authorize", service.Authorize)
	router.Post("/authorize", service.Authorize)

	router.Get("/token", service.Token)
	router.Post("/token", service.Token)
}

func (service *Routers) Authorize(w http.ResponseWriter, r *http.Request) {
	store, err := service.SessionManager.Start(r.Context(), w, r)
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
	store.Save()

	err = service.IServer.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (service *Routers) Token(w http.ResponseWriter, r *http.Request) {
	err := service.IServer.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
