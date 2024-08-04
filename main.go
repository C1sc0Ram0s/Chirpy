package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	database "github.com/C1sc0Ram0s/Chirpy/internal"
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

	router.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbConnection, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	chirps, err := dbConnection.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (id *idConfig) handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Id   int    `json:"id"`
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	id.id++
	chirps := returnVals{
		Id:   id.id,
		Body: cleaned,
	}

	// Creates new database connection and writes a chirp to disk
	dbConnection, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database connection failed")
	}
	dbConnection.CreateChirp(chirps.Body)
	respondWithJSON(w, 201, chirps)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, exists := badWords[lowerWord]; exists {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
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
