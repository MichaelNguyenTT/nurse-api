package routes

import (
	"net/http"
)

type NurseAPIRequests interface {
	GetSubjects(w http.ResponseWriter, r *http.Request)
}

// wrapping all handlers and middlewares
type NurseAPIWrapper struct {
	Handlers   NurseAPIRequests
	Middleware []Middlewares
	logEnabled bool
}

// layer the middlewares for logging/errors/auth etc
type Middlewares func(http.HandlerFunc) http.HandlerFunc

func (app *NurseAPIWrapper) GetSubjectsHandler(w http.ResponseWriter, r *http.Request) {
	// ...
	ctx := r.Context()

	var finalHandler = func(w http.ResponseWriter, r *http.Request) {
		app.Handlers.GetSubjects(w, r)
	}

	wrappedHandlers := app.applyMiddleware(finalHandler)
	wrappedHandlers(w, r.WithContext(ctx))
}
