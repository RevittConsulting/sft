package sft

import (
	"github.com/RevittConsulting/sft/sft/utils"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

// TODO: add some form of auth

type Handler struct {
	s *Service
}

func NewHandler(r chi.Router, s *Service) *Handler {
	h := &Handler{
		s: s,
	}
	h.SetupRoutes(r)
	return h
}

func (h *Handler) SetupRoutes(router chi.Router) {
	log.Println("setting up feature toggle routes.")

	router.Group(func(r chi.Router) {
		r.Route("/toggles", func(r chi.Router) {
			r.Get("/", h.GetAllToggles)
		})
	})
}

func (h *Handler) GetAllToggles(w http.ResponseWriter, r *http.Request) {
	toggles, err := h.s.GetAllToggles(r.Context())
	if err != nil {
		utils.WriteErr(w, err, http.StatusBadRequest)
	}

	utils.WriteJSON(w, toggles)
}
