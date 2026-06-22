package main

import (
    "os"
	"log"
	"net/http"
	"sync/atomic"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
    "github.com/jovinjoseph/chirpy-server/internal/database"
)

// apiConfig tracks stateful data across HTTP handlers.
// It is intended to be initialized once and shared between routes.
type apiConfig struct {
	fileserverHits atomic.Int32
    dbQueries *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"

    godotenv.Load()
    dbURL := os.Getenv("DB_URL")
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
	    log.Fatalf("Database connection failed: %v\n", err)
    }
    dbQueries := database.New(db)

	apiCfg := apiConfig{dbQueries: dbQueries}
	mux := http.NewServeMux()

	// The /app/ prefix is stripped so the file server can look up files
	fileServer := http.FileServer(http.Dir(filepathRoot))
	handler := http.StripPrefix("/app", fileServer)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
