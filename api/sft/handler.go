package sft

import (
	"encoding/json"
	"fmt"
	"github.com/RevittConsulting/sft/sft/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	go StartDashboard("6969")
	return h
}

func StartDashboard(port string) {

	// TODO: sort out path to the dashboard
	fs := http.FileServer(http.Dir("/Users/maxbb/github/revitt/sft/web/dashboard"))
	http.Handle("/", fs)

	log.Printf("Starting dashboard server on http://localhost:%s/", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Dashboard server failed to start: %v", err)
	}
}

func (h *Handler) SetupRoutes(router chi.Router) {
	log.Println("setting up feature toggle routes.")

	router.Group(func(r chi.Router) {
		r.Route("/toggles", func(r chi.Router) {
			r.Get("/", h.GetAllToggles)

			r.Post("/", h.CreateToggle)

			r.Patch("/disable/{toggle-id}", h.DisableFeature)
			r.Patch("/enable/{toggle_id}", h.EnableFeature)
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

func (h *Handler) CreateToggle(w http.ResponseWriter, r *http.Request) {
	log.Println("creating toggle")
	toggleDto := &ToggleDto{}

	err := json.NewDecoder(r.Body).Decode(toggleDto)
	if err != nil {
		utils.WriteErr(w, err, http.StatusBadRequest)
		return
	}

	newToggleId, err := h.s.CreateToggle(r.Context(), *toggleDto)

	utils.WriteJSON(w, newToggleId)
	return
}

func (h *Handler) DisableFeature(w http.ResponseWriter, r *http.Request) {
	log.Println("disabling feature")
	toggleId, err := uuid.Parse(chi.URLParam(r, "toggle-id"))
	if err != nil {
		message := fmt.Errorf("error parsing uuid: %w", err)
		utils.WriteErr(w, message, http.StatusBadRequest)
		return
	}

	err = h.s.DisableFeature(r.Context(), toggleId)
	if err != nil {
		utils.WriteErr(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, fmt.Sprintf("toggle %v set to disabled", toggleId))
	return
}

func (h *Handler) EnableFeature(w http.ResponseWriter, r *http.Request) {
	log.Println("enabling feature")
	toggleId, err := uuid.Parse(chi.URLParam(r, "toggle-id"))
	if err != nil {
		message := fmt.Errorf("error parsing uuid: %w", err)
		utils.WriteErr(w, message, http.StatusBadRequest)
		return
	}

	err = h.s.EnableFeature(r.Context(), toggleId)
	if err != nil {
		utils.WriteErr(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, fmt.Sprintf("toggle %v set to disabled", toggleId))
	return
}
