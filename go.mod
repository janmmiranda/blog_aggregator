module github.com/janmmiranda/blog_aggregator

go 1.22.2

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

require github.com/janmmiranda/blog_aggregator/internal/database v0.0.0

replace github.com/janmmiranda/blog_aggregator/internal/database => ./internal/database
