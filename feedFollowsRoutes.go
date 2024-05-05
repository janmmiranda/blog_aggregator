package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/janmmiranda/blog_aggregator/internal/database"
)

type FeedFollowResponse struct {
	Id        string    `json:"id"`
	FeedId    string    `json:"feed_id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) createFeedFollows(w http.ResponseWriter, req *http.Request, user database.User) {
	type params struct {
		FeedId string `json:"feed_id"`
	}

	feedFollowParam, err := decodeJson[params](req)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusBadRequest, "Error decoding feed follow create request")
		return
	}

	feedFollow, err := cfg.createFeedFollow(req.Context(), user.ID, feedFollowParam.FeedId)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Error creating new feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, convertFeedFollowResponse(feedFollow))
}

func (cfg *apiConfig) deleteFeedFollows(w http.ResponseWriter, req *http.Request) {
	feedFollowId := req.PathValue("feedFollowID")
	err := cfg.DB.DeleteFeedFollow(req.Context(), feedFollowId)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Unable to delete feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func (cfg *apiConfig) readFeedFollows(w http.ResponseWriter, req *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollows(req.Context(), user.ID)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusBadRequest, "Unable to fetch feed follows for user")
		return
	}
	feedFollowsResponse := make([]FeedFollowResponse, 0)
	for _, ff := range feedFollows {
		feedFollowsResponse = append(feedFollowsResponse, convertFeedFollowResponse(ff))
	}

	respondWithJSON(w, http.StatusOK, feedFollowsResponse)
}

func (cfg *apiConfig) createFeedFollow(ctx context.Context, userId, feedId string) (database.Feedfollow, error) {
	feedFollowId := uuid.NewString()
	currTime := time.Now().UTC()

	feedFollowCreate := database.CreateFeedFollowParams{
		ID:        feedFollowId,
		CreatedAt: currTime,
		UpdatedAt: currTime,
		UserID:    userId,
		FeedID:    feedId,
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(ctx, feedFollowCreate)
	if err != nil {
		return feedFollow, err
	}
	return feedFollow, nil
}

func convertFeedFollowResponse(feedFollowDb database.Feedfollow) FeedFollowResponse {
	return FeedFollowResponse{
		Id:        feedFollowDb.ID,
		FeedId:    feedFollowDb.FeedID,
		UserId:    feedFollowDb.UserID,
		CreatedAt: feedFollowDb.CreatedAt,
		UpdatedAt: feedFollowDb.UpdatedAt,
	}
}
