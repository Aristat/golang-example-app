package http

import (
	"net/http"
	"time"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/go-chi/chi/v5/middleware"
)

func Logger(l logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				l.Info("proto %s, path %s, lat %d, status %d, size %d, reqId %s", logger.Args(
					r.Proto,
					r.URL.Path,
					time.Since(t1),
					ww.Status(),
					ww.BytesWritten(),
					middleware.GetReqID(r.Context()),
				))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
