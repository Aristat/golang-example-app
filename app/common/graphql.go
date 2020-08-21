package common

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/99designs/gqlgen/graphql"
)

func SendGraphqlError(w http.ResponseWriter, code int, errors gqlerror.List) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	b, err := json.Marshal(&graphql.Response{Errors: errors})
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func SendGraphqlErrorf(w http.ResponseWriter, code int, format string, args ...interface{}) {
	SendGraphqlError(w, code, gqlerror.List{{Message: fmt.Sprintf(format, args...)}})
}
