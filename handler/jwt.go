package handler

import (
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateTokenPair(GUID uint, ip string) (map[string]string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS512)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["GUID"] = GUID
	claims["IP"] = ip
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	t, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, errors.New("Failed to create access token")
	}

	refreshToken := jwt.New(jwt.SigningMethodHS512)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["GUID"] = GUID
	rtClaims["IP"] = ip
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rt, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, errors.New("Failed to create refresh token")
	}

	rtBase64 := base64.StdEncoding.EncodeToString([]byte(rt))

	return map[string]string{
		"access_token":  t,
		"refresh_token": rtBase64,
	}, nil
}
