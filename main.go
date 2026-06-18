package main

import (
    "log"
    "net/http"
    "sync/atomic"
    "fmt"
)

// apiConfig holds shared, in-memory state for the server.
// Using a struct allows handler methods to access and mutate this state.
type apiConfig struct {
    // atomic.Int32 is a thread-safe integer - safe to read/write
    // across multiple concurrent goroutines (i.e., simultaneous HTTP requests).
    fileserverHits atomic.Int32
}

func main() {
    const filepathRoot = "."
	const port = "8080"

    // zero-value initialization is fine here - atomic.Int32 starts at 0 by default
    apiCfg := apiConfig{}

    // ServeMux is a request router - it maps URL patterns to handlers
    mux := http.NewServeMux()

    // /app/ route: wraps the file server with metrics middleware.
    // Reading right-to-left:
    //   1. http.Dir(filepathRoot)        - treats "." as the file system root
    //   2. http.FileServer(...)          - serves files from that root
    //   3. http.StripPrefix("/app", ...) - removes "/app" from the URL before
    //                                      passing it to the file server
    //   4. apiCfg.middlewareMetricsInc() - wraps the above, incrementing the
    //                                      hit counter on every request
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
    
    // /healthz: plain readiness check, no state needed, so a plain function works
    mux.HandleFunc("/healthz", handlerReadiness)

    // /metrics: returns the current hit count - needs apiCfg, so it's a method
    mux.HandleFunc("/metrics", apiCfg.handlerMetrics)

    // /reset: resets the hit count to 0 - also needs apiCfg
    mux.HandleFunc("/reset", apiCfg.handlerReset)

    srv := http.Server{
        Addr:    ":" + port, // listen on all interfaces at port 8080
        Handler: mux,        // route all requests through our mux
    }

    log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
    // ListenAndServe blocks forever, accepting and dispatching requests.
    // If it ever returns, it's because of a fatal error - log.Fatal handles that.
	log.Fatal(srv.ListenAndServe())
}

// middlewareMetricsInc wraps any http.Handler, incrementing the hit counter
// before passing the request along to the real handler (next).
// This is the middleware pattern: intercept -> do something -> continue.
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    // http.HandlerFunc converts a plain function into an http.Handler
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1) // atomically increment - safe under concurrency
        next.ServeHTTP(w, r)      // delegate to the wrapped handler
    })
}

// handlerMetrics responds with the current hit count as plain text.
// It's a method on *apiConfig so it can read cfg.fileserverHits.
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK) // 200 OK
    // .Load() atomically reads the current value of fileserverHits
    w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}