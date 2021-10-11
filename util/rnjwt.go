package util

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/model"
)

func CreateToken(user model.User) (string, error) {
	at := model.AuthTokenClaims{
		Id:       user.Id,
		Uuid:     NewUuid(),
		Role:     user.Role,
		Name:     user.Name,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	token, err := atoken.SignedString([]byte(config.JwtSecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func AuthToken(token string) (*jwt.Token, interface{}, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			ErrUnexpectedSigningMethod := errors.New("unexpected signing method")
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(config.JwtSecretKey), nil
	}
	claims := &model.AuthTokenClaims{}
	tok, err := jwt.ParseWithClaims(token, claims, keyFunc)
	return tok, *claims, err
}
