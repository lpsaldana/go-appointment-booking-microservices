package handlers

import "net/http"

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/login", h.login)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	println("hello")
}
