package sft

import (
	"encoding/json"
	"fmt"
	"github.com/RevittConsulting/logger"
	"github.com/RevittConsulting/sft/sft/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// TODO: add some form of auth

type Handler struct {
	s *Service
}

type Config struct {
	Buildpath string
	Port      string
}

func NewHandler(r chi.Router, s *Service, cfg *Config) *Handler {
	h := &Handler{
		s: s,
	}
	h.SetupRoutes(r)
	go StartDashboard(cfg)
	return h
}

func StartDashboard(cfg *Config) {

	//TODO: sort out path to the dashboard - this needs to be managed by a config when the dashboard is ready

	// Q for BOSSMAX: why is it we can't just grab the working directory and use that to find the buildpath (i.e. because the
	// build files will be within the library which is surely in the same place regardless of the application that is
	// implementing SFT)?

	// hardCodedPath := "/Users/maxbb/github/revitt/sft/web/dashboard/dist"

	buildPath := filepath.Join(cfg.Buildpath, "/sft/web/dashboard/dist")

	fmt.Println("buildpath is: ", buildPath)

	fs := http.FileServer(http.Dir(buildPath))

	// this handler function is necessary as we are serving a single page application (via React Router)
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

	logger.Log().Info("Starting dashboard server", zap.String("url", "http://localhost:"+cfg.Port))

	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatalf("Dashboard server failed to start: %v", err)
	}
}

func (h *Handler) SetupRoutes(router chi.Router) {
	logger.Log().Info("setting up feature toggle routes.")

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
	logger.Log().Info("Creating toggle")
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
	logger.Log().Info("toggling feature")
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
	logger.Log().Info("deleting toggle")
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
