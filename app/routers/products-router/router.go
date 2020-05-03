package products_router

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/aristat/golang-example-app/common"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"

	"github.com/go-chi/chi"

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
	chiRouter.Get("/products", router.GetProducts)
}

func (service *Router) GetProducts(w http.ResponseWriter, r *http.Request) {
	conn, d, err := grpc.GetConnGRPC(service.poolManager, common.SrvProducts)
	defer d()
	defer r.Body.Close()

	if err != nil {
		service.logger.Printf("[ERROR] %s", err.Error())
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

	// Connect to NATS
	nc, err := nats.Connect(service.cfg.NatsURL)
	if err != nil {
		service.logger.Error(err.Error())
	}
	defer nc.Close()

	sc, err := stan.Connect("test-cluster", "stan-pub", stan.NatsConn(nc))
	message := "Hello"
	service.logger.Printf("[NATS] send %s", message)
	err = sc.Publish(service.cfg.Subject, []byte(message))
	if err != nil {
		service.logger.Printf("[ERROR] %s", err.Error())
	}

	// Close connection
	sc.Close()

	e.Encode(productOut)
}
