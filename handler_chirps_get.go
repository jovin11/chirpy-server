package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jovinjoseph/chirpy-server/internal/database"
)

func (cfg *apiConfig) handlerChirpsGetByID(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	response := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, response)

}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {

	var dbChirps []database.Chirp
	var err error

	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		userID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
		dbChirps, err = cfg.dbQueries.GetChirpsByAuthorID(r.Context(), userID)

	} else {
		dbChirps, err = cfg.dbQueries.GetChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
		return
	}

	chirps := make([]Chirp, 0, len(dbChirps))

	for _, dbChirp := range dbChirps {
		chirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
		chirps = append(chirps, chirp)
	}
	respondWithJSON(w, http.StatusOK, chirps)

}
