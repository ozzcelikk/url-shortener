package main

import (
	"context"
	"encoding/json"
	"github.com/ozzcelikk/url-shortener/src/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func httpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		entry := log.New(os.Stdout, "[ENTRY] ", 0)
		entry.Printf("Time: ", start)
		entry.Printf("URL: ", uri)
		entry.Printf("Method: ", method)

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		exit := log.New(os.Stdout, "[EXIT] ", 0)
		exit.Printf("Duration: ", duration)
	})
}
func rootRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		services.GetHandler(w, r)
	case http.MethodPost:
		w.Header().Set("content-type", "application/json")
		services.CreateHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func listRouteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	switch r.Method {
	case http.MethodGet:
		res := services.ListHandler()
		json.NewEncoder(w).Encode(res)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func serviceCheckHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	services.Init()

	mux := http.NewServeMux()

	serviceHandler := http.HandlerFunc(serviceCheckHandler)
	mux.Handle("/service", httpMiddleware(serviceHandler))

	rootHandler := http.HandlerFunc(rootRouteHandler)
	mux.Handle("/", httpMiddleware(rootHandler))

	lstRouteHandler := http.HandlerFunc(listRouteHandler)
	mux.Handle("/list", httpMiddleware(lstRouteHandler))

	srv := &http.Server{
		Handler:      mux,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
