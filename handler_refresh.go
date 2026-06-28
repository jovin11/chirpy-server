package main

import (
	"github.com/jovinjoseph/chirpy-server/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	if refreshToken.RevokedAt.Valid || refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is revoked or expired", nil)
		return
	}

	token, err = auth.MakeJWT(
		refreshToken.UserID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create access JWT", err)
		return
	}

	res := response{
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, res)
}
