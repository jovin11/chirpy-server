package main

import (
    "log"
    "net/http"
)

func main() {

    mux := http.NewServeMux()
    server := http.Server{
        Addr: ":8080", 
        Handler: mux,
    }
    err := server.ListenAndServe()

    if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
