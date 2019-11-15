package common

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/gqlerror"
)

func SendGraphqlError(w http.ResponseWriter, code int, errors ...*gqlerror.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	b, err := json.Marshal(&graphql.Response{Errors: errors})
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func SendGraphqlErrorf(w http.ResponseWriter, code int, format string, args ...interface{}) {
	SendGraphqlError(w, code, &gqlerror.Error{Message: fmt.Sprintf(format, args...)})
}
