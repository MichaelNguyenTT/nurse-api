package routes

import (
	"net/http"
)

// use this middleware to apply more layers of checks
func (api *NurseAPIWrapper) applyMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	if api.logEnabled {
		handler = loggingMiddleware(handler)
	}

	// auth....?

	for _, middleware := range api.Middleware {
		handler = middleware(handler)
	}

	return handler
}
