package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type APIHandler struct {
	repo   *Repository
	logger *slog.Logger
}

func NewRepo(db *gorm.DB) *APIHandler {
	return &APIHandler{
		repo: NewRepository(db),
	}
}

func (a *APIHandler) List(w http.ResponseWriter, r *http.Request) {
	serv, err := a.repo.List()
	if err != nil {
		a.logger.Error("database access failure", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(serv.GetAll())
	if err != nil {
		a.logger.Error("JSON response error", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (a *APIHandler) Create(w http.ResponseWriter, r *http.Request) {
	req := &ServiceRequest{}

	decodeJSON(w, r.Body, req)

	newService := req.ToDB()
	newService.ID = uuid.New()

	_, err := a.repo.Create(newService)
	if err != nil {
		a.logger.Error("failed to connect to database", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *APIHandler) Read(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error("invalid url params", "msg", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	serv, err := a.repo.Read(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			a.logger.Error("no database record", "msg", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.logger.Error("failed to connect to database", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := serv.ToClient()
	encodeJSON(w, res)
}

func (a *APIHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error("invalid url params", "msg", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqModel := &ServiceRequest{}
	decodeJSON(w, r.Body, reqModel)

	serve := reqModel.ToDB()
	serve.ID = id

	db, err := a.repo.Update(serve)
	if err != nil {
		a.logger.Error("failed to connect to database", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if db == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (a *APIHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error("invalid url params", "msg", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := a.repo.Delete(id)
	if err != nil {
		a.logger.Error("failed to connect to database", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if db == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}