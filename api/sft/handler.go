package sft

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/RevittConsulting/logger"
	"github.com/RevittConsulting/sft/sft/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// TODO: add some form of auth

type Handler struct {
	s *Service
}

type Config struct {
	Buildpath       string
	Port            string
	ApplicationName string
	ApiAddress      string
}

func NewHandler(r chi.Router, s *Service, cfg *Config) *Handler {
	h := &Handler{
		s: s,
	}
	h.SetupRoutes(r)
	go StartDashboard(cfg)
	return h
}

//go:embed build_artifacts/dist/*
var web embed.FS
var uiFS fs.FS

func StartDashboard(cfg *Config) {

	var err error
	uiFS, err = fs.Sub(web, "build_artifacts/dist")
	if err != nil {
		log.Fatal("failed to get ui fs", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/*", handleStatic)
	r.Get("/api/config.json", func(w http.ResponseWriter, r *http.Request) {
		handleApiConfig(w, r, cfg)
	})

	log.Println("starting server...")
	if err := http.ListenAndServe(":6969", r); err != nil {
		log.Println("server failed:", err)
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	path := filepath.Clean(r.URL.Path)
	if path == "/" {
		path = "index.html"
	} else {
		path = strings.TrimPrefix(path, "/")
	}

	// Attempt to open the requested file
	file, err := uiFS.Open(path)
	if err != nil {
		// Fallback to index.html for unknown routes
		log.Println("file", path, "not found:", err)
		file, err = uiFS.Open("index.html")
		if err != nil {
			log.Println("file index.html cannot be read:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		path = "index.html"
	}
	defer file.Close()

	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	if strings.HasPrefix(path, "static/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}

	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	n, _ := io.Copy(w, file)
	log.Println("file", path, "copied", n, "bytes")
}

func handleApiConfig(w http.ResponseWriter, r *http.Request, cfg *Config) {
	log.Println("In the handleAPI function")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg.ApiAddress)
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
