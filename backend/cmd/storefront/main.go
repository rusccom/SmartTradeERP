package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"smarterp/backend/internal/platform/storefront"
	"smarterp/backend/internal/shared/app"
	"smarterp/backend/internal/shared/config"
)

func main() {
	port := flag.Int("preview-port", 0, "Serve a local MamaPack draft with sample products on the loopback interface")
	flag.Parse()
	if *port != 0 {
		runPreview(*port)
		return
	}
	runStorefront()
}

func runStorefront() {
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

func runPreview(port int) {
	handler, err := storefront.LocalPreview()
	if err != nil {
		log.Fatal(err)
	}
	addr := "127.0.0.1:" + strconv.Itoa(port)
	server := http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("MamaPack preview: http://%s/", addr)
	log.Fatal(server.ListenAndServe())
}
