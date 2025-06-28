package main

import (
	"log"
	"net/http"
	"pet-mgt/backend/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.LoadCfg()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	srv := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	log.Printf("listening on port %s", cfg.Port)
	srv.ListenAndServe()
}
