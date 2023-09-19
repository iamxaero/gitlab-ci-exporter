package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"example.com/gitlab-ci-exporter/config"
	"example.com/gitlab-ci-exporter/controller"
	"github.com/cloudflare/cfssl/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Level = log.LevelDebug
	log.Infof("gitlab-ci-exporter have been started")
	// Config
	cfg := config.New()
	ctrl := controller.New(cfg)
	// Router
	h2s := &http2.Server{}
	handler := http.NewServeMux()
	// Prometheus register metrics
	ctrl.PromRegister()
	// Handlers
	handler.HandleFunc("/webhook", ctrl.Webhook)
	handler.HandleFunc("/", ctrl.Health)
	handler.HandleFunc("/health", ctrl.Health)
	handler.Handle("/metrics", promhttp.Handler())
	// Set option for server
	listen := ":8080"
	server := &http.Server{
		Addr:         listen,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 180 * time.Second,
		IdleTimeout:  240 * time.Second,
		Handler:      h2c.NewHandler(handler, h2s),
	}
	// Start http server
	go func() {
		log.Infof("Running server at %v", listen)
		log.Fatal(server.ListenAndServe())
	}()
	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
