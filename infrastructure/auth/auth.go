package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// AccessDetail ...
type AccessDetail struct {
	TokenUUID string
	UserID    string
}

// TokenDetail ...
type TokenDetail struct {
	AccessToken        string `json:"access"`
	RefreshToken       string `json:"refresh"`
	TokenUUID          string `json:"token_uuid"`
	RefreshUUID        string `json:"refresh_uuid"`
	AccessTokenExpire  int64  `json:"token_expire"`
	RefreshTokenExpire int64  `json:"refresh_expire"`
}

// AccessData ...
type AccessData struct {
	redisDB *redis.Client
}

// AuthenticationInterface ...
type AuthenticationInterface interface {
	CreateToken(string, *TokenDetail) ([]string, error)
	FetchToken(string) (string, error)
	DeleteAccessToken(*AccessDetail) error
	DeleteRefreshToken(string) error
}

// NewAccessData ...
func NewAccessData(redisDB *redis.Client) *AccessData {
	return &AccessData{redisDB: redisDB}
}

var _ AuthenticationInterface = &AccessData{}

var ctx = context.Background()

// CreateToken ...
func (access *AccessData) CreateToken(userID string, tokenDetail *TokenDetail) ([]string, error) {
	isLoggedIn, _ := access.redisDB.Keys(ctx, fmt.Sprintf("*++%s", userID)).Result()
	accessExpire := time.Unix(tokenDetail.AccessTokenExpire, 0)
	refreshExpire := time.Unix(tokenDetail.RefreshTokenExpire, 0)
	timeNow := time.Now()

	accessCreated, err := access.redisDB.Set(ctx, tokenDetail.TokenUUID, userID, accessExpire.Sub(timeNow)).Result()
	if err != nil {
		return isLoggedIn, nil
	}

	refreshCreated, err := access.redisDB.Set(ctx, tokenDetail.RefreshUUID, userID, refreshExpire.Sub(timeNow)).Result()
	if err != nil {
		return isLoggedIn, nil
	}

	if accessCreated == "0" || refreshCreated == "0" {
		return isLoggedIn, errors.New("general_error")
	}
	return isLoggedIn, nil
}

// FetchToken ...
func (access *AccessData) FetchToken(tokenUUID string) (string, error) {
	userID, err := access.redisDB.Get(ctx, tokenUUID).Result()
	if err != nil {
		return "", errors.New("general_error")
	}
	return userID, nil
}

// DeleteAccessToken ...
func (access *AccessData) DeleteAccessToken(auth *AccessDetail) error {
	refreshUUID := fmt.Sprintf("%s++%s", auth.TokenUUID, auth.UserID)
	_, err := access.redisDB.Del(ctx, auth.TokenUUID).Result()

	if err != nil {
		return errors.New("general_error")
	}

	_, err = access.redisDB.Del(ctx, refreshUUID).Result()
	if err != nil {
		return nil
	}
	return nil
}

// DeleteRefreshToken ...
func (access *AccessData) DeleteRefreshToken(refreshUUID string) error {
	_, err := access.redisDB.Del(ctx, refreshUUID).Result()
	if err != nil {
		return errors.New("general_error")
	}
	return nil
}
