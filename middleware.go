package main

import (
	"log"
	"net/http"

	"github.com/janmmiranda/blog_aggregator/internal/database"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		apikey, err := GetToken(req.Header, "ApiKey")
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusUnauthorized, "Couldn't find api key")
			return
		}
		user, err := cfg.DB.ReadUserByAPIKey(req.Context(), apikey)
		if err != nil {
			log.Println(err.Error())
			respondWithError(w, http.StatusNotFound, "Couldn't get user")
			return
		}

		handler(w, req, user)
	}

}
