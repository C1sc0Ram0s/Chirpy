package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}
type idConfig struct {
	id int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	idCfg := idConfig{
		id: 0,
	}

	router := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/*", fsHandler)

	router.HandleFunc("GET /api/healthz", handlerReadiness)
	router.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	router.HandleFunc("POST /api/chirps", idCfg.handlerChirpsValidate)
	router.HandleFunc("GET /api/chirps", handlerGetChirps)
	router.HandleFunc("GET /api/chirps/{chirpID}", handlerGetChirpID)

	router.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
