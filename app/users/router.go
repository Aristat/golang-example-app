package users

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"

	"github.com/aristat/golang-example-app/generated/resources/proto/products"

	"google.golang.org/grpc"

	"github.com/jinzhu/gorm"

	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/aristat/golang-example-app/app/logger"
	"gopkg.in/oauth2.v3/server"

	"github.com/go-session/session"

	"github.com/aristat/golang-example-app/common"
	"github.com/go-chi/chi"
)

type Router struct {
	ctx            context.Context
	sessionManager *session.Manager
	template       *template.Template
	logger         logger.Logger
	db             *gorm.DB
	server         *server.Server
	repo           *repo.Repo
}

func (router *Router) Run(chiRouter *chi.Mux) {
	chiRouter.Get("/login", router.GetLogin)
	chiRouter.Post("/login", router.PostLogin)

	chiRouter.Get("/auth", router.Auth)
	chiRouter.Post("/auth", router.Auth)

	chiRouter.Get("/user", router.User)

	chiRouter.Get("/products", router.GetProducts)
}

func (service *Router) GetProducts(w http.ResponseWriter, r *http.Request) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts,
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),
		)))
	opts = append(opts, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
		otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer()),
	)))

	conn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer conn.Close()
	c := products.NewProductsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	productOut, err := c.ListProduct(ctx, &products.ListProductIn{Id: 1})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	e := json.NewEncoder(w)
	e.Encode(productOut)
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

	currentUser, err := service.repo.Users.FindByEmail(email)
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
