package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/janmmiranda/blog_aggregator/internal/database"
)

type PostResponse struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      string    `json:"feed_id"`
}
type Pagination struct {
	NextOffset     int `json:"next_offset"`
	PreviousOffset int `json:"previous_offset"`
}

type PostsResponse struct {
	Pagination Pagination     `json:"pagination"`
	Posts      []PostResponse `json:"posts"`
}

func (cfg *apiConfig) readPostsByUser(w http.ResponseWriter, req *http.Request, user database.User) {
	limit, err := strconv.Atoi(req.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		log.Println("Unable to fetch query limit")
		limit = 10
	}

	offset, err := strconv.Atoi(req.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	nextOffset := offset + limit
	prevOffset := offset - limit
	if prevOffset < 0 {
		prevOffset = 0
	}

	postsDB, err := cfg.DB.GetPostsByUser(req.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusNoContent, "Unable to fetch posts for user")
		return
	}
	postsResponse := make([]PostResponse, 0)
	for _, post := range postsDB {
		postsResponse = append(postsResponse, convertDBPostToResponse(post))
	}

	respondWithJSON(w, http.StatusFound, PostsResponse{
		Pagination: Pagination{
			NextOffset:     nextOffset,
			PreviousOffset: prevOffset,
		},
		Posts: postsResponse,
	})
}

func convertDBPostToResponse(postDB database.Post) PostResponse {
	return PostResponse{
		ID:          postDB.ID,
		CreatedAt:   postDB.CreatedAt,
		UpdatedAt:   postDB.UpdatedAt,
		Title:       postDB.Title.String,
		Url:         postDB.Url,
		Description: postDB.Description.String,
		PublishedAt: postDB.PublishedAt.Time,
		FeedID:      postDB.FeedID,
	}
}
