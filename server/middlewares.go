package server

import (
	"context"
	"net/http"

	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/routes"
)

// HealthcheckMiddleware is used to set the database in the current context.
func HealthcheckMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI == "/status" {
				w.WriteHeader(http.StatusOK)

				return
			}

			next.ServeHTTP(w, r.WithContext(r.Context()))
		}

		return http.HandlerFunc(fn)
	}
}

// ValidateQueryParametersMiddleware validates that the incoming request has the proper query parameters
// if not it is descarted.
// NOTE: also it stores the values into a context if the are found.
func ValidateQueryParametersMiddleware(qrps []routes.RequiredQueryParameter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Validate that all query parameters that are required are passed correctly
			for i := range qrps {
				// if the current request doesnt have the query parameter then return an
				// status no found
				// TODO: we can return difference thing specifying the error but not sure so far.
				v := r.URL.Query().Get(string(qrps[i]))
				if v == "" {
					w.WriteHeader(http.StatusNotFound)

					return
				}

				// save the current query parameter and its value
				ctx = context.WithValue(ctx, qrps[i], v)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
