package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/sync/errgroup"

	"bnbl.io/csp_exporter/collector"
	"bnbl.io/csp_exporter/csp"
)

type config struct {
	CollectorBindAddr string `envconfig:"COLLECTOR_BIND_ADDR" default:":80"`
	PromBindAddr      string `envconfig:"PROM_BIND_ADDR" default:":9477"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("could not read configuration from environment: %v", err)
	}

	coll, err := collector.NewCollector()
	if err != nil {
		log.Fatalf("could not create collector: %v", err)
	}

	cspHandler := http.NewServeMux()
	cspHandler.HandleFunc("/report/csp/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Header.Get("content-type") != "application/csp-report" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		app := strings.TrimPrefix(r.URL.EscapedPath(), "/report/csp/")
		if app == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		report, err := csp.ReadReport(r.Body)
		if err == io.EOF {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		coll.Collect(app, report)
	})

	g, ctx := errgroup.WithContext(context.Background())

	// start metrics server
	g.Go(func() error {
		h := http.NewServeMux()
		h.Handle("/metrics", coll.Handler())
		srv := &http.Server{Addr: cfg.PromBindAddr, Handler: h}
		go func() {
			<-ctx.Done()
			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("shutdown metrics server: %v", err)
			}
		}()
		return fmt.Errorf("metrics server: %v", srv.ListenAndServe())
	})

	// start CSP collector server
	g.Go(func() error {
		srv := &http.Server{Addr: cfg.CollectorBindAddr, Handler: cspHandler}
		go func() {
			<-ctx.Done()
			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("shutdown collector server: %v", err)
			}
		}()
		return fmt.Errorf("collector server: %v", srv.ListenAndServe())
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
