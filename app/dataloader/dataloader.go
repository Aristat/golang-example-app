//go:generate go run github.com/vektah/dataloaden ProductItemLoader int []*github.com/aristat/golang-example-app/generated/domain.ProductItem

package dataloader

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aristat/golang-example-app/app/db/domain"
)

type ctxKeyType struct{ name string }

var ctxKey = ctxKeyType{"userDataLoader"}

type loaders struct {
	ProductItemsByProduct *ProductItemLoader
}

func LoaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ldrs := loaders{}

		// set this to zero what happens without data loading
		wait := 250 * time.Microsecond

		ldrs.ProductItemsByProduct = &ProductItemLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(keys []int) ([][]*domain.ProductItem, []error) {
				fmt.Println("ProductItems Start Fetch")
				var keySql []string
				for _, key := range keys {
					keySql = append(keySql, strconv.Itoa(key))
				}

				fmt.Printf("SELECT * FROM product_items WHERE product_id IN (%s)\n", strings.Join(keySql, ","))
				time.Sleep(5 * time.Millisecond)

				productItems := make([][]*domain.ProductItem, len(keys))
				errors := make([]error, len(keys))
				for i := range keys {
					productItems[i] = []*domain.ProductItem{
						{ID: 1, Name: "item " + strconv.Itoa(rand.Int()%20+20)},
						{ID: 2, Name: "item " + strconv.Itoa(rand.Int()%20+20)},
					}
				}

				return productItems, errors
			},
		}

		dlCtx := context.WithValue(r.Context(), ctxKey, ldrs)
		next.ServeHTTP(w, r.WithContext(dlCtx))
	})
}

func CtxLoaders(ctx context.Context) loaders {
	return ctx.Value(ctxKey).(loaders)
}
