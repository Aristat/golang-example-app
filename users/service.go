package users

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"gopkg.in/oauth2.v3/server"

	"github.com/jinzhu/gorm"

	"github.com/go-chi/chi"

	"fmt"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"github.com/aristat/golang-gin-oauth2-example-app/common"
	"github.com/go-session/session"
)

var (
	userNotFound = errors.New("10002 user not found")
)

type H map[string]interface{}

type Service struct {
	Template       *template.Template
	SessionManager *session.Manager
	DB             *gorm.DB
	Log            logger.Logger
	*server.Server
}

func Run(router *chi.Mux, service *Service) {
	router.Get("/login", service.GetLogin)
	router.Post("/login", service.PostLogin)

	router.Get("/auth", service.Auth)
	router.Post("/auth", service.Auth)

	router.Get("/user", service.User)
}

func (service *Service) GetLogin(w http.ResponseWriter, r *http.Request) {
	_, err := service.SessionManager.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.Template.ExecuteTemplate(w, "users/login", H{})
}

func (service *Service) PostLogin(w http.ResponseWriter, r *http.Request) {
	store, err := service.SessionManager.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	currentUser, err := FindByEmail(service.DB, email)
	if err != nil {
		service.Template.ExecuteTemplate(w, "users/login", H{"errors": []string{userNotFound.Error()}})
		return
	}

	if common.CheckPasswordHash(password, currentUser.EncryptedPassword) == false {
		service.Template.ExecuteTemplate(w, "users/login", H{"errors": []string{userNotFound.Error()}})
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

func (service *Service) Auth(w http.ResponseWriter, r *http.Request) {
	store, err := service.SessionManager.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		service.Log.Error("User not found")

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)

		return
	}

	service.Template.ExecuteTemplate(w, "users/auth", H{})
}

func (service *Service) User(w http.ResponseWriter, r *http.Request) {
	ti, err := service.Server.ValidationBearerToken(r)
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
