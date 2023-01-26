package net

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

type service interface {
	Login(ctx context.Context, params *model.LoginInfo) (*model.Authentication, error)
	Verify(ctx context.Context, params *model.Authentication) (bool, error)
	Refresh(ctx context.Context, params *model.Authentication) (*model.Authentication, error)
}

type handler struct {
	service service
}

func New(service service) *handler {
	return &handler{
		service: service}
}

func (h *handler) Verify(w http.ResponseWriter, r *http.Request) {
	value := r.Header.Get("Authorization")

	token := strings.TrimPrefix(value, "Bearer ") // header Authorization bare
	params := model.Authentication{
		Token: token,
	}

	auth, err := h.service.Verify(r.Context(), &params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(auth); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params model.LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	auth, err := h.service.Login(r.Context(), &params)
	if err != nil || auth == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(auth); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("token", auth.Token)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	value := r.Header.Get("Authorization")

	token := strings.TrimPrefix(value, "Bearer ") // header Authorization bare
	params := model.Authentication{
		Token: token,
	}

	auth, err := h.service.Refresh(r.Context(), &params)
	if err != nil {
		// TODO return unAuth
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(auth); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("token", auth.Token)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) Handler() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		})
		r.Post("/verify", h.Verify)
		r.Post("/access_token", h.Login)
		r.Post("/refresh", h.Refresh)
	})

	return r
}
