package server

import (
	"fmt"
	"net/http"

	"github.com/TurniXXD/go-rss/internal/auth"
	"github.com/TurniXXD/go-rss/internal/database"
	"github.com/TurniXXD/go-rss/utils"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			utils.JSONErrorResponse(w, 401, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			utils.JSONErrorResponse(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
