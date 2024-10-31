package service

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	e "nms/api/service/errors"
)

func decodeJSON(w http.ResponseWriter, r io.ReadCloser, val any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r).Decode(val)
	if err != nil {
		slog.Error("failed to decode json", "msg", err)
		e.ResponseInternalErr(w)
		return
	}
}

func encodeJSON(w http.ResponseWriter, val any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(val)
	if err != nil {
		slog.Error("Failed to encode", "msg", err)
		e.ResponseInternalErr(w)
		return
	}
}
