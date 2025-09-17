package main

import (
	"fmt"
	"net/http"

	"example.com/internal/api"
	"example.com/internal/config"
	"example.com/internal/recorder"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	rec := recorder.NewRecorder()
	server := api.NewServer(cfg, rec)
	router := server.NewRouter()

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("Server running on http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}

}
