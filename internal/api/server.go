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

	r.Get("/ping", handlePing)
	r.Get("/status", handleStatus)
	r.Get("/config", s.getConfig)
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

// POST /config
func (s *Server) updateConfig(w http.ResponseWriter, r *http.Request) {
	var newCfg config.Config
	if err := json.NewDecoder(r.Body).Decode(&newCfg); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// overwrite live config
	*s.cfg = newCfg

	// persist to file
	if err := config.Save(s.cfg); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}
