package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/janmmiranda/blog_aggregator/internal/database"
)

type FeednFollowResponse struct {
	Feed       FeedResponse       `json:"feed"`
	FeedFollow FeedFollowResponse `json:"feed_follow"`
}

type FeedResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserId    string    `json:"user_id"`
}

func (cfg *apiConfig) createFeed(w http.ResponseWriter, req *http.Request, user database.User) {
	type params struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	reqFeed, err := decodeJson[params](req)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusBadRequest, "Error decoding feed create request")
		return
	}

	feedId := uuid.NewString()
	currTime := time.Now().UTC()

	feedCreate := database.CreateFeedParams{
		ID:        feedId,
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Name:      reqFeed.Name,
		Url:       reqFeed.Url,
		UserID:    user.ID,
	}

	feed, err := cfg.DB.CreateFeed(req.Context(), feedCreate)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Error creating new feed")
		return
	}

	feedFollowDb, err := cfg.createFeedFollow(req.Context(), user.ID, feed.ID)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Error creating new feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, FeednFollowResponse{
		Feed: FeedResponse{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserId:    feed.UserID,
		},
		FeedFollow: convertFeedFollowResponse(feedFollowDb),
	})
}

func (cfg *apiConfig) readFeeds(w http.ResponseWriter, req *http.Request) {
	feeds, err := cfg.DB.ReadFeeds(req.Context())
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusNotFound, "Couldn't find feeds")
		return
	}

	respondWithJSON(w, http.StatusFound, convertFeedsForResponse(feeds))
}

func convertFeedsForResponse(feedsDb []database.Feed) []FeedResponse {
	feedsResponse := make([]FeedResponse, 0)
	for _, feed := range feedsDb {
		feedsResponse = append(feedsResponse, FeedResponse{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserId:    feed.UserID,
		})
	}
	return feedsResponse
}
