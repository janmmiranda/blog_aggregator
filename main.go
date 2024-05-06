package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/janmmiranda/blog_aggregator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	CTX *context.Context
	DB  *database.Queries
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	port := os.Getenv("PORT")
	conn := os.Getenv("CONN")

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal("Error connecting to postgres")
	}
	dbQueries := database.New(db)

	apiConfig := apiConfig{
		CTX: &ctx,
		DB:  dbQueries,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/readiness", handlerReadiness)
	mux.HandleFunc("GET /v1/error", handlerError)

	mux.HandleFunc("POST /v1/users", apiConfig.createUser)
	mux.HandleFunc("GET /v1/users", apiConfig.middlewareAuth(apiConfig.readUser))

	mux.HandleFunc("POST /v1/feeds", apiConfig.middlewareAuth(apiConfig.createFeed))
	mux.HandleFunc("GET /v1/feeds", apiConfig.readFeeds)

	mux.HandleFunc("POST /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.createFeedFollows))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiConfig.deleteFeedFollows)
	mux.HandleFunc("GET /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.readFeedFollows))

	mux.HandleFunc("GET /v1/posts", apiConfig.middlewareAuth(apiConfig.readPostsByUser))

	corsMux := middlewareLog(middlewareCors(mux))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	go func() {
		<-ctx.Done() // Wait for cancellation signal
		log.Println("Shutting down server...")

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Attempt to shut down the server.
		if err := server.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Fatalf("Error shutting down server: %v", err)
		}
	}()

	const collectionConcurrency = 10
	const collectionInterval = 10 * time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Starting server on port %s...\n", port)
	return server.ListenAndServe()
}
