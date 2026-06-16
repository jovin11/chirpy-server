package main

import (
    "log"
    "net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
     w.Header().Add("Content-Type", "text/plain; charset=utf-8")
     w.WriteHeader(http.StatusOK)
     w.Write([]byte("OK"))
}

func main() {

    mux := http.NewServeMux()

    mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
    mux.HandleFunc("/healthz", handlerReadiness)
    

    server := http.Server{
        Addr: ":8080", 
        Handler: mux,
    }
    err := server.ListenAndServe()

    if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
