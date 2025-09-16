package api

import (
	"encoding/json"
	"net/http"

	"example.com/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//GET
	r.Get("/ping", handlePing)
	r.Get("/status", handleStatus)
	r.Get("/config", s.getConfig)

	//POST
	r.Post("/config", s.updateConfig)

	return r
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (s *Server) getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.cfg)
}

func (s *Server) updateConfig(w http.ResponseWriter, r *http.Request) {
	var updates config.ConfigUpdate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if updates.Server.Port != nil {
		s.cfg.Server.Port = *updates.Server.Port
	}
	if updates.Output.Directory != nil {
		s.cfg.Output.Directory = *updates.Output.Directory
	}
	if updates.Output.FPS != nil {
		s.cfg.Output.FPS = *updates.Output.FPS
	}
	if updates.Output.StartVideoRecordingHotkey != nil {
		s.cfg.Output.StartVideoRecordingHotkey = *updates.Output.StartVideoRecordingHotkey
	}
	if updates.Output.CaptureScreenShotHotkey != nil {
		s.cfg.Output.CaptureScreenShotHotkey = *updates.Output.CaptureScreenShotHotkey
	}

	// persist to file
	if err := config.Save(s.cfg); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}
