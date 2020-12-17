package products_router

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/nats-io/stan.go"

	"github.com/aristat/golang-example-app/app/common"
	"github.com/go-chi/chi"
	"github.com/nats-io/nats.go"

	"github.com/aristat/golang-example-app/app/grpc"
	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/generated/resources/proto/products"
)

type Router struct {
	ctx         context.Context
	cfg         *Config
	template    *template.Template
	logger      logger.Logger
	poolManager *grpc.PoolManager
}

func (router *Router) Run(chiRouter chi.Router) {
	chiRouter.Get("/products_grpc", router.GetProductsGrpc)
	chiRouter.Get("/products_nats", router.GetProductsNats)
	chiRouter.Get("/products_slowly", router.GetProductsSlowly)
}

func (router *Router) GetProductsGrpc(w http.ResponseWriter, r *http.Request) {
	conn, d, err := grpc.GetConnGRPC(router.poolManager, common.SrvProducts)
	defer d()
	defer r.Body.Close()

	if err != nil {
		router.logger.Printf("[ERROR] %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := products.NewProductsClient(conn)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	e := json.NewEncoder(w)

	productOut, err := c.ListProduct(ctx, &products.ListProductIn{Id: 1})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		e.Encode("{}")
		return
	}

	e.Encode(productOut)
}

func (router *Router) GetProductsNats(w http.ResponseWriter, r *http.Request) {
	// Connect to NATS
	nc, err := nats.Connect(router.cfg.NatsURL)
	if err != nil {
		router.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer nc.Close()

	sc, err := stan.Connect("test-cluster", "stan-pub", stan.NatsConn(nc))
	if err != nil {
		router.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Close connection
	defer sc.Close()

	message := "Hello"
	router.logger.Printf("[NATS] send %s", message)
	err = sc.Publish(router.cfg.Subject, []byte(message))

	if err != nil {
		router.logger.Printf("[ERROR] %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(w)
	e.Encode("done")
}

func (router *Router) GetProductsSlowly(w http.ResponseWriter, r *http.Request) {
	router.logger.Info("Start sleep")
	time.Sleep(time.Second * 10)
	router.logger.Info("Stop sleep")
	e := json.NewEncoder(w)
	e.Encode("")
}
