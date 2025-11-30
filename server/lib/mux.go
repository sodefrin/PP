package lib

import (
	"log/slog"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type ServeMux struct {
	mux *http.ServeMux
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		mux: http.NewServeMux(),
	}
}

func (m *ServeMux) Handle(pattern string, handler http.Handler) {
	m.mux.Handle(pattern, handler)
}

func (m *ServeMux) HandleFunc(pattern string, handler HandlerFunc) {
	m.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			slog.ErrorContext(r.Context(), "internal server error", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}
