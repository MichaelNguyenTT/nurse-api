package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type FinalHandler func(http.HandlerFunc) http.HandlerFunc

type MuxOpts struct {
	HostPath    string
	Router      *mux.Router
	Middlewares []Middlewares
}

// handle all the paths
func HandlerConfigurations(n NurseAPIRequests, m MuxOpts) http.Handler {
	mux := m.Router

	wrapper := NurseAPIWrapper{
		Handlers:   n,
		Middleware: m.Middlewares,
		logEnabled: true,
	}

	s := mux.PathPrefix(m.HostPath + "/v1/api").Subrouter()
	s.Path("/subjects").Methods("GET").HandlerFunc(wrapper.GetSubjectsHandler)

	return mux
}

func NewHandler(api NurseAPIRequests, r *mux.Router) http.Handler {
	return HandlerConfigurations(api, MuxOpts{
		Router: r,
	})
}
