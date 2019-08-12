package users

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"
	"gopkg.in/oauth2.v3/server"

	"github.com/go-session/session"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
	"github.com/go-chi/chi"
)

type Router struct {
	ctx            context.Context
	sessionManager *session.Manager
	template       *template.Template
	logger         logger.Logger
	db             *gorm.DB
	server         *server.Server
}

func (router *Router) Run(chiRouter *chi.Mux) {
	chiRouter.Get("/login", router.GetLogin)
	chiRouter.Post("/login", router.PostLogin)

	chiRouter.Get("/auth", router.Auth)
	chiRouter.Post("/auth", router.Auth)

	chiRouter.Get("/user", router.User)
}

func (service *Router) GetLogin(w http.ResponseWriter, r *http.Request) {
	_, err := service.sessionManager.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.template.ExecuteTemplate(w, "users/login", H{})
}

func (service *Router) PostLogin(w http.ResponseWriter, r *http.Request) {
	store, err := service.sessionManager.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	currentUser, err := FindByEmail(service.db, email)
	if err != nil {
		service.template.ExecuteTemplate(w, "users/login", H{"errors": []string{userNotFound.Error()}})
		return
	}

	if common.CheckPasswordHash(password, currentUser.EncryptedPassword) == false {
		service.template.ExecuteTemplate(w, "users/login", H{"errors": []string{userNotFound.Error()}})
		return
	}

	if currentUser != nil {
		store.Set("LoggedInUserID", fmt.Sprintf("%d", currentUser.ID))
		store.Save()

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}

	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusFound)
}

func (service *Router) Auth(w http.ResponseWriter, r *http.Request) {
	store, err := service.sessionManager.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		service.logger.Error("User not found")

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)

		return
	}

	service.template.ExecuteTemplate(w, "users/auth", H{})
}

func (service *Router) User(w http.ResponseWriter, r *http.Request) {
	ti, err := service.server.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	status := http.StatusOK

	w.WriteHeader(status)
	userData := userData{ID: ti.GetUserID(), Scope: ti.GetScope()}
	err = json.NewEncoder(w).Encode(userData)
}
