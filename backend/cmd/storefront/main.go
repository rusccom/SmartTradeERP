package main

import (
	"context"
	"log"

	"smarterp/backend/internal/shared/app"
	"smarterp/backend/internal/shared/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	server, cleanup, err := app.BuildStorefront(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	log.Printf("storefront listening on %s", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
