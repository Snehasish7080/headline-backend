package jwtclaim

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

func CreateJwtToken(id string) (string, error) {

	claims := UserClaim{
		id,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(240 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	var jwtKey = []byte("my_secret_key")
	jsonWebToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jsonWebToken.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractId(tokenStr string) (string, bool) {
	hmacSecretString := "my_secret_key"
	hmacSecret := []byte(hmacSecretString)

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return "", false
	}
	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return claims.Id, ok
	} else {
		return "", false
	}
}
