package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/karimOCB/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform bool
}

type User struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
}

func main() {
	port := "8080"
	cfg := apiConfig{}
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	cfg.platform = os.Getenv("PLATFORM") == "dev" 
		
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("could establish connection to database: %v", err)
	}

	cfg.dbQueries = database.New(db)
	
	mux := http.NewServeMux() 

	server := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	filepathRoot := http.Dir(".")
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(filepathRoot))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	/*TODO: pass errors to respondWithError and use http named constant to improve readibility*/	
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't listen to the server: %v", err)
	}
}