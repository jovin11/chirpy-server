package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
    type requestParams struct {
		Email string `json:"email"`
	}
	type UserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	var params requestParams
    err := json.NewDecoder(r.Body).Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Failed to decode parameters", err)
        return
    }

    // call cfg.db.CreateUser(...)
	user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

    // map database user to response struct
	response := UserResponse {
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	// respond with JSON and created status
	respondWithJSON(w, http.StatusCreated, response)
    
}