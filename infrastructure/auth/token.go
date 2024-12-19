package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/majid-cj/go-chat-server/util"

	"github.com/golang-jwt/jwt/v5"
)

// TokenInterface ...
type TokenInterface interface {
	CreateJWTToken(string, string, string) (*TokenDetail, error)
	ExtractJWTTokenMetadata(*http.Request, bool) (*AccessDetail, error)
}

// Token ...
type Token struct{}

var _ TokenInterface = &Token{}

// NewToken ...
func NewToken() *Token {
	return &Token{}
}

// CreateJWTToken ...
func (token *Token) CreateJWTToken(userId, profileId, uniqueId string) (*TokenDetail, error) {
	tokenDetail := &TokenDetail{}
	tokenDetail.AccessTokenExpire = time.Now().Add(time.Hour * 24).Unix()
	tokenDetail.TokenUUID = util.ULID()

	tokenDetail.RefreshTokenExpire = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenDetail.RefreshUUID = fmt.Sprintf("%s++%s", tokenDetail.TokenUUID, userId)

	var err error
	accessTokenClaim := jwt.MapClaims{}
	accessTokenClaim["authorization"] = true
	accessTokenClaim["access_uuid"] = tokenDetail.TokenUUID
	accessTokenClaim["user_id"] = userId
	accessTokenClaim["profile_id"] = profileId
	accessTokenClaim["unique_id"] = uniqueId
	accessTokenClaim["exp"] = tokenDetail.AccessTokenExpire
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)

	tokenDetail.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, errors.New("general_error")
	}

	refreshTokenClaim := jwt.MapClaims{}
	refreshTokenClaim["refresh_uuid"] = tokenDetail.RefreshUUID
	refreshTokenClaim["user_id"] = userId
	refreshTokenClaim["profile_id"] = profileId
	refreshTokenClaim["unique_id"] = uniqueId
	refreshTokenClaim["exp"] = tokenDetail.RefreshTokenExpire
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaim)

	tokenDetail.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, errors.New("general_error")
	}
	return tokenDetail, nil
}

// ExtractJWTTokenMetadata ...
func (token *Token) ExtractJWTTokenMetadata(request *http.Request, checkToken bool) (*AccessDetail, error) {
	_token, err := VerifyToken(request, checkToken)
	if err != nil {
		return nil, errors.New("general_error")
	}
	claims, ok := _token.Claims.(jwt.MapClaims)
	if ok && _token.Valid && checkToken {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("general_error")
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("general_error")
		}
		return &AccessDetail{
			TokenUUID: accessUUID,
			UserID:    userID,
		}, nil
	}
	if ok && !checkToken {
		accessUUID, _ := claims["access_uuid"].(string)
		userID, _ := claims["user_id"].(string)
		return &AccessDetail{
			TokenUUID: accessUUID,
			UserID:    userID,
		}, nil
	}
	return nil, errors.New("general_error")
}

// ExtractTokenClaims ...
func ExtractTokenClaims(request *http.Request, key string) string {
	token, err := VerifyToken(request, true)
	if err != nil {
		return ""
	}
	if !token.Valid {
		return ""
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	value, ok := claims[key].(string)
	if !ok {
		return ""
	}
	return value
}

// ExtractURLTokenClaims ...
func ExtractURLTokenClaims(request *http.Request, key string) string {
	token, err := VerifyURLToken(request)
	if err != nil {
		return ""
	}
	if !token.Valid {
		return ""
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	value, ok := claims[key].(string)
	if !ok {
		return ""
	}
	return value
}

// TokenValid ...
func TokenValid(request *http.Request) error {
	token, err := VerifyToken(request, true)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		// invalid token
		return errors.New("general_error")
	}
	return nil
}

// URLTokenValid ...
func URLTokenValid(request *http.Request) bool {
	token, err := VerifyURLToken(request)
	if err != nil {
		return false
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		// invalid token
		return false
	}
	return true
}

// VerifyToken ...
func VerifyToken(request *http.Request, checkToken bool) (*jwt.Token, error) {
	_token := ExtractToken(request)
	token, err := jwt.Parse(_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if checkToken && err != nil {
		// error parsing token
		return nil, errors.New("general_error")
	}
	return token, nil
}

// VerifyURLToken ...
func VerifyURLToken(request *http.Request) (*jwt.Token, error) {
	jwtToken := ExtractURLToken(request)
	value, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, errors.New("general_error")
	}
	return value, nil
}

// ExtractToken ...
func ExtractToken(request *http.Request) string {
	bearer := request.Header.Get("Authorization")
	token := strings.Split(bearer, " ")
	if token[0] != "Bearer" {
		return ""
	}
	if len(token) == 2 {
		return token[1]
	}
	return ""
}

// ExtractURLToken ...
func ExtractURLToken(request *http.Request) string {
	return request.URL.Query().Get("access")
}
