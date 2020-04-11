package models

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims for jwt token
type Claims struct {
	ID uint64
	jwt.StandardClaims
}

func IssueToken(ID uint64, expiration time.Time) *Claims {
	return &Claims{
		ID: ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}
}

// ParseIfValidJWT parses a JWT token if it is valid
func (c *Claims) ParseIfValidJWT(token string) error {
	_, err := jwt.ParseWithClaims(
		token,
		c,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCOUNT_JWT_KEY")), nil
		},
	)
	return err
}
