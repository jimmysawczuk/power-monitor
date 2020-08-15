package respond

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type errResponse struct {
	Error     string `json:"error"`
	ErrorCode string `json:"code,omitempty"`
	Status    int    `json:"status,omitempty"`
}

// WithError is a shortcut for WithCodedError(log, w, httpStatus, "", err).
func WithError(log logrus.FieldLogger, w http.ResponseWriter, r *http.Request, httpStatus int, err error) {
	WithCodedError(log, w, r, httpStatus, "", err)
}

// WithCodedError writes the provided error to the ResponseWriter, as well as the HTTP status code.
// An enum-style code (i.e. INVALID_TOKEN) may also be provided.
func WithCodedError(log logrus.FieldLogger, w http.ResponseWriter, r *http.Request, httpStatus int, code string, err error) {
	serr := http.StatusText(httpStatus)
	if err != nil {
		serr = err.Error()
	}

	if ct := w.Header().Get("Content-Type"); ct == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(errResponse{
		Error:     serr,
		Status:    httpStatus,
		ErrorCode: code,
	})
}

// WithSuccess wraps the provided response with a success field and sets the provided HTTP response status.
func WithSuccess(log logrus.FieldLogger, w http.ResponseWriter, r *http.Request, httpStatus int, v interface{}) {
	if ct := w.Header().Get("Content-Type"); ct == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(v)
}
