package auth

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// DBAuth ...
type DBAuth struct {
	Auth AuthenticationInterface
	DB   *redis.Client
}

// NewDBAuth ...
func NewDBAuth() *DBAuth {
	HOST := os.Getenv("AUTH_HOST")
	PORT := os.Getenv("AUTH_PORT")
	PASSWORD := os.Getenv("AUTH_PASSWORD")
	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", HOST, PORT),
		Password: PASSWORD,
		DB:       0,
	})

	return &DBAuth{
		Auth: NewAccessData(db),
		DB:   db,
	}
}
