package apiserver

import (
	"fmt"
	"net/http"

	"github.com/roi-05/kvstore-raft/internal/store"
)

type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/set":
		key := r.URL.Query().Get("key")
		val := r.URL.Query().Get("val")
		if key == "" || val == "" {
			http.Error(w, "missing key or val", http.StatusBadRequest)
			return
		}
		h.store.Set(key, val)
		fmt.Fprint(w, "OK")

	case "/get":
		key := r.URL.Query().Get("key")
		val, ok := h.store.Get(key)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		fmt.Fprint(w, val)

	default:
		http.NotFound(w, r)
	}
}
