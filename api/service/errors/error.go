package errors

import (
	"encoding/json"
	"net/http"
)

const (
	ErrInternalServer = "internal server error"
	ErrNotFound       = "not found"
	ErrBadRequest     = "bad request"
	ErrInvalidID      = "invalid id format"
)

type ServerError struct {
	Error string `json:"error"`
}

func runError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ServerError{Error: msg})
}

func ResponseInternalErr(w http.ResponseWriter) {
	runError(w, ErrInternalServer, http.StatusInternalServerError)
}

func ResponseNotFoundErr(w http.ResponseWriter) {
	runError(w, ErrNotFound, http.StatusNotFound)
}

func ResponseBadRequestErr(w http.ResponseWriter) {
	runError(w, ErrBadRequest, http.StatusBadRequest)
}

func ResponseInvalidID(w http.ResponseWriter) {
	runError(w, ErrInvalidID, http.StatusBadRequest)
}
