package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jovinjoseph/chirpy-server/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	polkaKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}
	if polkaKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}

	_, err = cfg.dbQueries.UpdateUserChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
