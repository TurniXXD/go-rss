package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TurniXXD/go-rss/internal/database"
	"github.com/TurniXXD/go-rss/utils"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feedId"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.JSONErrorResponse(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		utils.JSONErrorResponse(w, 400, fmt.Sprintf("Couldn't create feed follow: %v", err))
		return
	}

	utils.JSONResponse(w, 201, mapDbFeedFollow(feed))
}

func (cfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		utils.JSONErrorResponse(w, 400, fmt.Sprintf("Couldn't get feed follows: %v", err))
		return
	}

	utils.JSONResponse(w, 200, mapDbFeedFollowsToFeedFollows(feedFollows))
}

func (cfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID := chi.URLParam(r, "feedFollowID")
	feedFollowIDParsed, err := uuid.Parse(feedFollowID)
	if err != nil {
		utils.JSONErrorResponse(w, 400, fmt.Sprintf("Error parsing feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowIDParsed,
		UserID: user.ID,
	})

	if err != nil {
		utils.JSONErrorResponse(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
	}

	utils.JSONResponse(w, 200, struct{}{})
}
