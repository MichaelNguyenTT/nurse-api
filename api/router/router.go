package router

import (
	"nms/api/service"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		mux := service.NewRepo(db)
		r.Get("/service", mux.List)
		r.Post("/service", mux.Create)
		r.Get("/service/{id}", mux.Read)
		r.Put("/service/{id}", mux.Update)
		r.Delete("/service/{id}", mux.Delete)
	})

	return r
}
