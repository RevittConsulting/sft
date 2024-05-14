package sft

import (
	"encoding/json"
	"fmt"
	"github.com/RevittConsulting/sft/sft/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
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

	//TODO: sort out path to the dashboard - this needs to be managed by a config when the dashboard is ready
	buildPath := "/Users/maxbb/github/revitt/sft/web/dashboard/dist"

	fs := http.FileServer(http.Dir(buildPath))

	// this handler function is designed to allow us to serve a react app with react router
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the requested file exists
		path := buildPath + r.URL.Path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// If the file does not exist, serve index.html
			http.ServeFile(w, r, buildPath+"/index.html")
		} else {
			// Otherwise, serve the file
			fs.ServeHTTP(w, r)
		}
	})

	log.Printf("Starting dashboard server on http://hello:%s/", port)
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

			r.Patch("/{toggle-id}", h.ToggleFeature)

			r.Delete("/{toggle-id}", h.DeleteToggle)
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

func (h *Handler) ToggleFeature(w http.ResponseWriter, r *http.Request) {
	log.Println("toggling feature")
	toggleId, err := uuid.Parse(chi.URLParam(r, "toggle-id"))
	if err != nil {
		message := fmt.Errorf("error parsing uuid: %w", err)
		utils.WriteErr(w, message, http.StatusBadRequest)
		return
	}

	err = h.s.ToggleFeature(r.Context(), toggleId)
	if err != nil {
		utils.WriteErr(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, fmt.Sprintf("toggle %v set to disabled", toggleId))
	return
}

func (h *Handler) DeleteToggle(w http.ResponseWriter, r *http.Request) {
	log.Println("deleting toggle")
	toggleId, err := uuid.Parse(chi.URLParam(r, "toggle-id"))
	if err != nil {
		message := fmt.Errorf("error parsing uuid: %w", err)
		utils.WriteErr(w, message, http.StatusBadRequest)
		return
	}

	err = h.s.DeleteToggle(r.Context(), toggleId)
	if err != nil {
		utils.WriteErr(w, err, http.StatusInternalServerError)
		return
	}

}
