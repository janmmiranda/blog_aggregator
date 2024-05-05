package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, response{
		Status: "ok",
	})
}

func handlerError(w http.ResponseWriter, req *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
