package products_router

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/aristat/golang-example-app/app/grpc"
	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/common"
	"github.com/aristat/golang-example-app/generated/resources/proto/products"
)

type Router struct {
	ctx         context.Context
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
	e.Encode(productOut)
}
