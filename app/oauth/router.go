package oauth

import (
	"context"
	"net/http"
	"net/url"

	"gopkg.in/oauth2.v3/server"

	"github.com/go-chi/chi"
)

type Router struct {
	ctx     context.Context
	Server  *server.Server
	Service *Service
}

func (router *Router) Run(chiRouter *chi.Mux) {
	chiRouter.Get("/authorize", router.Authorize)
	chiRouter.Post("/authorize", router.Authorize)

	chiRouter.Get("/token", router.Token)
	chiRouter.Post("/token", router.Token)
}

func (router *Router) Authorize(w http.ResponseWriter, r *http.Request) {
	store, err := router.Service.SessionManager.Start(r.Context(), w, r)
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
