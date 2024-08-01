package main

import (
	"fmt"
	"log"
	"net/http"
)

type Handler struct{}

func (Handler) ServeHTTP(http.ResponseWriter, *http.Request) {}

// Add a handler which displays status
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// ApiConfig holds our application state
type ApiConfig struct {
	fileserverHits int
}

// Middleware function to increment request count
func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits++

		// Call the next handler
		next.ServeHTTP(w, req)
	})
}

// Add a handler which displays metrics
func (cfg *ApiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	</html>
	`, cfg.fileserverHits)))
}

// Add a handler which resets the metrics
func (cfg *ApiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	var apiCfg = &ApiConfig{}

	router := http.NewServeMux()
	router.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	router.HandleFunc("GET /api/healthz", handlerReadiness)
	router.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	router.HandleFunc("/api/reset", apiCfg.handlerReset)

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
