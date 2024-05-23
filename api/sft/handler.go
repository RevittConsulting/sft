package sft

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/RevittConsulting/logger"
	"github.com/RevittConsulting/sft/sft/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io/fs"
	"log"
	"net/http"
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

func StartDashboard(cfg *Config) {

	//// Create a sub-filesystem from the embedded files, pointing directly at the dist folder
	//distFS, err := fs.Sub(web, "build_artifacts/dist")
	//if err != nil {
	//	log.Fatalf("Failed to initialize embedded filesystem: %v", err)
	//}
	//
	//// File server to serve static assets
	//staticHandler := http.FileServer(http.FS(distFS))
	//
	//// Custom handler to deal with SPA routing
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	// Log the incoming request URL
	//	log.Printf("Requested URL: %s", r.URL.Path)
	//
	//	// Normalize the path to avoid bypassing the handler logic
	//	path := strings.TrimPrefix(r.URL.Path, "/")
	//	if path == "" || path == "dashboard" || path == "create" {
	//		// Log decision to serve index.html for SPA paths
	//		log.Printf("Serving index.html for SPA route: %s", path)
	//
	//		// Directly serve index.html for SPA routes and root
	//		path = "index.html"
	//
	//		// Check if index.html exists in the embedded file system
	//		file, err := distFS.Open(path)
	//		if err != nil {
	//			// Log if there's an error opening index.html
	//			log.Printf("Error opening index.html: %v", err)
	//			http.Error(w, "Internal server error", http.StatusInternalServerError)
	//			return
	//		}
	//		fileStat, err := file.Stat()
	//		if err != nil {
	//			log.Printf("Error getting stats for index.html: %v", err)
	//			http.Error(w, "Internal server error", http.StatusInternalServerError)
	//			return
	//		}
	//		// Log file details
	//		log.Printf("Serving file: %s, Size: %d", fileStat.Name(), fileStat.Size())
	//		http.StripPrefix("/", staticHandler).ServeHTTP(w, r)
	//	}
	//
	//	// Attempt to open the file
	//	_, err := distFS.Open(path)
	//	if err != nil {
	//		// Log file open error
	//		log.Printf("Error opening file '%s': %v", path, err)
	//		// Fallback to serving index.html if there's an error
	//		http.ServeFile(w, r, "build_artifacts/dist/index.html")
	//		return
	//	}
	//
	//	// Log static file serving
	//	log.Printf("Serving static file: %s", path)
	//	staticHandler.ServeHTTP(w, r)
	//})
	//
	//// Start the server
	//logger.Log().Info("Starting dashboard server", zap.String("url", "http://localhost:"+cfg.Port))
	//err = http.ListenAndServe(":"+cfg.Port, nil)
	//if err != nil {
	//	log.Fatalf("Dashboard server failed to start: %v", err)
	//}

	entries, err := fs.ReadDir(web, "build_artifacts/dist")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Contents of dist:")
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}

	fsys, err := fs.Sub(web, "build_artifacts/dist")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", spaHandler(http.FS(fsys), cfg.Buildpath))

	log.Println("Serving on http://localhost:6969")
	err = http.ListenAndServe(":6969", nil)
	if err != nil {
		log.Fatal(err)
	}

	//// SPA METHOD THAT ONLY SERVES THE ROOT PAGE!!!
	//// Create a sub-filesystem from the embedded files, pointing directly at the dist folder
	//distFS, err := fs.Sub(web, "build_artifacts/dist")
	//if err != nil {
	//	log.Fatalf("Failed to initialize embedded filesystem: %v", err)
	//}
	//
	//// File server to serve static assets
	//staticHandler := http.FileServer(http.FS(distFS))
	//
	//// Custom handler to deal with SPA routing
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	// Log the incoming request URL
	//	log.Printf("Requested URL: %s", r.URL.Path)
	//
	//	// Normalize the path to avoid bypassing the handler logic
	//	path := strings.TrimPrefix(r.URL.Path, "/")
	//	if path == "" || path == "dashboard" || path == "create" {
	//		// Log decision to serve index.html for SPA paths
	//		log.Printf("Serving index.html for SPA route: %s", path)
	//
	//		// Directly serve index.html for SPA routes and root
	//		path = "index.html"
	//
	//		// Check if index.html exists in the embedded file system
	//		file, err := distFS.Open(path)
	//		if err != nil {
	//			// Log if there's an error opening index.html
	//			log.Printf("Error opening index.html: %v", err)
	//			http.Error(w, "Internal server error", http.StatusInternalServerError)
	//			return
	//		}
	//		fileStat, err := file.Stat()
	//		if err != nil {
	//			log.Printf("Error getting stats for index.html: %v", err)
	//			http.Error(w, "Internal server error", http.StatusInternalServerError)
	//			return
	//		}
	//		// Log file details
	//		log.Printf("Serving file: %s, Size: %d", fileStat.Name(), fileStat.Size())
	//		http.StripPrefix("/", staticHandler).ServeHTTP(w, r)
	//	}
	//
	//	// Attempt to open the file
	//	_, err := distFS.Open(path)
	//	if err != nil {
	//		// Log file open error
	//		log.Printf("Error opening file '%s': %v", path, err)
	//		// Fallback to serving index.html if there's an error
	//		http.ServeFile(w, r, "build_artifacts/dist/index.html")
	//		return
	//	}
	//
	//	// Log static file serving
	//	log.Printf("Serving static file: %s", path)
	//	staticHandler.ServeHTTP(w, r)
	//})
	//
	//// Start the server
	//logger.Log().Info("Starting dashboard server", zap.String("url", "http://localhost:"+cfg.Port))
	//err = http.ListenAndServe(":"+cfg.Port, nil)
	//if err != nil {
	//	log.Fatalf("Dashboard server failed to start: %v", err)
	//}

	//// EMBED METHOD, SHOWS A SINGLE INDEX.HTML
	//
	//dist, _ := fs.Sub(web, "build_artifacts/dist")
	//
	//http.Handle("/", http.FileServer(http.FS(dist)))
	//
	//logger.Log().Info("Starting dashboard server", zap.String("url", "http://localhost:"+cfg.Port))
	//
	//err := http.ListenAndServe(":"+cfg.Port, nil)
	//if err != nil {
	//	log.Fatalf("Dashboard server failed to start: %v", err)
	//}

	// ORIGINAL METHOD
	//// hardCodedPath := "/Users/maxbb/github/revitt/sft/web/dashboard/dist"
	//
	//buildPath := filepath.Join(cfg.Buildpath, "/sft/web/dashboard/dist")
	//
	//fmt.Println("buildpath is: ", buildPath)
	//
	//fs := http.FileServer(http.Dir(buildPath))
	//
	//// this handler function is necessary as we are serving a single page application (via React Router)
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	// Check if the requested file exists
	//	path := buildPath + r.URL.Path
	//	if _, err := os.Stat(path); os.IsNotExist(err) {
	//		// If the file does not exist, serve index.html
	//		http.ServeFile(w, r, buildPath+"/index.html")
	//	} else {
	//		// Otherwise, serve the file
	//		fs.ServeHTTP(w, r)
	//	}
	//})
	//
	//logger.Log().Info("Starting dashboard server", zap.String("url", "http://localhost:"+cfg.Port))
	//
	//err := http.ListenAndServe(":"+cfg.Port, nil)
	//if err != nil {
	//	log.Fatalf("Dashboard server failed to start: %v", err)
	//}

}

func spaHandler(fsys http.FileSystem, buildPath string) http.Handler {
	fmt.Println("hello from the function")
	//fileServer := http.FileServer(fsys) // Create the file server for the embedded filesystem

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		fmt.Println("243", path)

		// Check if the file exists in the embedded filesystem
		_, err := fsys.Open(strings.TrimPrefix(path, "/"))
		if err != nil {
			// If the file does not exist, serve the index.html
			http.ServeFile(w, r, "/index.html")
			return
		}

		// Otherwise, serve the file
		http.FileServer(fsys).ServeHTTP(w, r)
	})
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
