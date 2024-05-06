module github.com/janmmiranda/blog_aggregator

go 1.22.2

require (
	github.com/google/uuid v1.6.0
	github.com/janmmiranda/blog_aggregator/internal/database v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	golang.org/x/net v0.24.0
)

require golang.org/x/text v0.14.0 // indirect

replace github.com/janmmiranda/blog_aggregator/internal/database => ./internal/database
