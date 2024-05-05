package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/janmmiranda/blog_aggregator/internal/database"
)

type UserResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Apikey    string    `json:"api_key"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {
	type params struct {
		Name string `json:"name"`
	}
	reqName, err := decodeJson[params](req)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Error decoding user create request")
		return
	}

	userId := uuid.NewString()
	currTime := time.Now().UTC()
	apikey, err := generateAPIKey()
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Error creating user's apikey")
		return
	}

	userReq := database.CreateUserParams{
		ID:        userId,
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Name:      reqName.Name,
		Apikey:    apikey,
	}

	user, err := cfg.DB.CreateUser(req.Context(), userReq)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Error creating new user")
		return
	}

	respondWithJSON(w, http.StatusCreated, UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Apikey:    user.Apikey,
	})
}

func (cfg *apiConfig) readUser(w http.ResponseWriter, req *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Apikey:    user.Apikey,
	})
}

func generateAPIKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	keyHex := hex.EncodeToString(key)
	return keyHex, nil
}
