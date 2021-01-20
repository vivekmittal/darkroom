package middleware

import (
	"github.com/gojek/darkroom/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
)

func Recovery() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logger.SugaredWithRequest(r)

			defer func() {
				err := recover()
				if err != nil {
					l.Errorf("Recovered from panic: %v", err)

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

