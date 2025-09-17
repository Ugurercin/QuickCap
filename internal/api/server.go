package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"example.com/internal/config"
	"example.com/internal/recorder"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	cfg *config.Config
	rec *recorder.Recorder
}

func NewServer(cfg *config.Config, recorder *recorder.Recorder) *Server {
	return &Server{cfg: cfg, rec: recorder}
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

	r.Post("/record/start", s.startRecording)
	r.Post("/record/stop", s.stopRecording)
	r.Get("/record/status", s.recordStatus)

	return r
}

func (s *Server) startRecording(w http.ResponseWriter, r *http.Request) {
	cfg := s.cfg.Output

	if err := os.MkdirAll(cfg.Directory, 0755); err != nil {
		http.Error(w, "failed to create output directory", http.StatusInternalServerError)
		return
	}

	file := filepath.Join(cfg.Directory,
		fmt.Sprintf("rec_%d.mp4", time.Now().Unix()))

	if err := s.rec.Start(file, *s.cfg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":        "recording",
		"file":          file,
		"mode":          string(cfg.Mode),
		"fps":           cfg.FPS,
		"position":      string(cfg.Position),
		"webcam_width":  cfg.WebcamW,
		"webcam_height": cfg.WebcamH,
	})
}

func (s *Server) stopRecording(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	file, err := s.rec.Stop(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "stopped",
		"file":   file,
	})
}

func (s *Server) recordStatus(w http.ResponseWriter, r *http.Request) {
	running, file, elapsed := s.rec.Status()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"running": running,
		"file":    file,
		"elapsed": elapsed,
	})
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

	if updates.Server != nil && updates.Server.Port != nil {
		s.cfg.Server.Port = *updates.Server.Port
	}
	if updates.Output != nil {
		if updates.Output.Directory != nil {
			s.cfg.Output.Directory = *updates.Output.Directory
		}
		if updates.Output.FPS != nil {
			s.cfg.Output.FPS = *updates.Output.FPS
		}
		/*	if updates.Output.StartVideoRecordingHotkey != nil {
			s.cfg.Output.StartVideoRecordingHotkey = *updates.Output.StartVideoRecordingHotkey
		}*/
		/*if updates.Output.CaptureScreenShotHotkey != nil {
			s.cfg.Output.CaptureScreenShotHotkey = *updates.Output.CaptureScreenShotHotkey
		}*/
		if updates.Output.Mode != nil {
			s.cfg.Output.Mode = *updates.Output.Mode
		}
		if updates.Output.Position != nil {
			s.cfg.Output.Position = *updates.Output.Position
		}
		if updates.Output.WebcamW != nil {
			s.cfg.Output.WebcamW = *updates.Output.WebcamW
		}
		if updates.Output.WebcamH != nil {
			s.cfg.Output.WebcamH = *updates.Output.WebcamH
		}
	}

	config.ValidateEnums(s.cfg)

	if err := config.Save(s.cfg); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.cfg)
}
