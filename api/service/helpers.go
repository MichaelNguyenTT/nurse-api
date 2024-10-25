package service

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func decodeJSON(w http.ResponseWriter, r io.ReadCloser, val any) {
	err := json.NewDecoder(r).Decode(val)
	if err != nil {
		slog.Error("failed to decode json", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func encodeJSON(w http.ResponseWriter, val any) {
	err := json.NewEncoder(w).Encode(val)
	if err != nil {
		slog.Error("Failed to encode", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
