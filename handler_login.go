package main

import (
	"encoding/json"
	"github.com/jovinjoseph/chirpy-server/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		UserResponse
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	var params requestParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode parameters", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create access JWT", err)
		return
	}

	res := response{
		UserResponse: UserResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        token,
		RefreshToken: "",
	}

	respondWithJSON(w, http.StatusOK, res)

}
