package routegroup

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger is an example middleware that logs the request method and URL and the time it took to process the request.
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		slog.Info("http request", "method", r.Method, "url", r.URL.Path, "took", time.Since(start))
	})
}

// Recover is an example middleware that recovers from panics,
// logs the panic, and returns a HTTP 500 (Internal Server Error)
// response status if possible.
func Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic", "error", err, "method", r.Method, "url", r.URL.Path)

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}
