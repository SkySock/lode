package v1

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Hello struct {
	l *slog.Logger
}

func NewHello(l *slog.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Info("Hello")

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Hello %s", d)
}
