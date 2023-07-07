package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TurniXXD/go-rss/internal/database"
	"github.com/TurniXXD/go-rss/utils"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.JSONErrorResponse(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		utils.JSONErrorResponse(w, 400, fmt.Sprintf("Couldn't create feed: %v", err))
		return
	}

	utils.JSONResponse(w, 201, mapDbFeedToFeed(feed))
}

func (cfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		utils.JSONErrorResponse(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}

	utils.JSONResponse(w, 200, mapDbFeedsToFeeds(feeds))
}
