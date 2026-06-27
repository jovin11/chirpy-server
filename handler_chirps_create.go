package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jovinjoseph/chirpy-server/internal/auth"
	"github.com/jovinjoseph/chirpy-server/internal/database"
	"net/http"
	"strings"
	"time"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned_body := getCleanedBody(params.Body, badWords)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: userID,
	}

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp", err)
		return
	}

	response := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, response)

}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
