package common

import (
	"BTaskServer/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("a_secret")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

type ClaimsRoot struct {
	Info string
	jwt.StandardClaims
}

func ReleseToken(user model.User) (string, error) {
	//expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			//ExpiresAt: expirationTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer:   "author",
			Subject:  "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ReleseRootToken(info string) (string, error) {
	expirationTime := time.Now().Add(1 * 24 * time.Hour)
	claims := &ClaimsRoot{
		Info: info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "author",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
