package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"nms/internal/db"
	d "nms/internal/db"
	"strconv"
	"sync"
)

var (
	StatusInternalServerError = http.StatusInternalServerError
)

type SubjectHandler struct {
	storage *d.SubjectCollection
	l       *slog.Logger
	mu      sync.RWMutex
}

func NewHandler(db *db.SubjectCollection) *SubjectHandler {
	return &SubjectHandler{
		storage: db,
	}
}

func (p *SubjectHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	subs, err := p.storage.All()
	if err != nil {
		p.l.Error("internal database error", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var subResponder []SubjectRequests
	for _, data := range subs {
		subResponder = append(subResponder, processDBToResponder(data))
	}

	EncodeJSONWriter(w, subResponder)
}

func (p *SubjectHandler) GetSubject(w http.ResponseWriter, r *http.Request, rawSubId SubjectID) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	patientId, err := strconv.Atoi(string(rawSubId))
	if err != nil {
		p.l.Error("couldn't strconv patient id", "msg", patientId)
		w.WriteHeader(StatusInternalServerError)
		return
	}

	patient, err := p.storage.ByID(patientId)
	if err != nil {
		p.l.Error("failed to retrieve patient information", "msg", err)
		if errors.Is(err, d.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(StatusInternalServerError)
		}
		return
	}

	processor := processDBToResponder(patient)
	EncodeJSONWriter(w, processor)
}

// TODO: finish the handles later
func (p *SubjectHandler) PostPatient(w http.ResponseWriter, r *http.Request) {

}

func (p *SubjectHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {

}

// from db to client transformer
func processDBToResponder(p d.SubjectDB) SubjectRequests {
	return SubjectRequests{
		ID:        p.ID,
		Name:      p.Name,
		Category:  p.Category,
		Priority:  p.Priority,
		CreatedAt: p.CreatedAt,
	}
}

func DecodeJSONWriter(w http.ResponseWriter, r *http.Request, data any) {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error("encountered error decoding request body", "msg", err, "key", data)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func EncodeJSONWriter(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&data)
	if err != nil {
		slog.Error("encountered error encoding response body", "msg", err, "key", data)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
