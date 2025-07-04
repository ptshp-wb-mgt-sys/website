// Package main contains the main entrypoint for the application
package main

import (
	"log"
	"net/http"
	"pet-mgt/backend/internal/config"
	"pet-mgt/backend/internal/routes"
)

func main() {
	cfg, err := config.LoadCfg()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	r := routes.SetupRouter(cfg)

	srv := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	log.Printf("listening on port %s", cfg.Port)
	srv.ListenAndServe()
}
