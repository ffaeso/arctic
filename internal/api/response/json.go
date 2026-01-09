package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func WriteJSON[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func ReadJSON[T any](r *http.Request) (T, error) {
	var v T
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

type ErrorResponse struct {
	StatusCode int    `json:"code,omitempty"`
	Status     string `json:"status,omitempty"`
	ErrorMsg   string `json:"message,omitempty"`
	Path       string `json:"path,omitempty"`
	Timestamp  string `json:"timestamp"`
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, message string) {
	e := ErrorResponse{
		ErrorMsg:   message,
		StatusCode: status,
		Status:     http.StatusText(status),
		Path:       r.URL.Path,
		Timestamp:  time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}
	if err := WriteJSON(w, status, e); err != nil {
		slog.ErrorContext(r.Context(), "failed to write error response", "error", err)
	}
}
