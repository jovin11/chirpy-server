package main

import (
    "net/http"
    "errors"
)

// handlerReset resets the statistics tracked by the server.
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
        respondWithError(w, http.StatusForbidden, "Forbidden access", errors.New("Reset is only allowed in dev environment"))
        return
    }

    err := cfg.dbQueries.DeleteUsers(r.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "failed to reset the db", err)
        return
    }

    cfg.fileserverHits.Store(0)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}