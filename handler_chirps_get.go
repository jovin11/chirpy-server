package main

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	type ChirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirps, err := cfg.dbQueries.GetChirps(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to receive chirps", err)
		return
	}

	responses := make([]ChirpResponse, 0, len(chirps))

	for _, chirp := range chirps {

		response := ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		responses = append(responses, response)
	}
	respondWithJSON(w, http.StatusOK, responses)

}
