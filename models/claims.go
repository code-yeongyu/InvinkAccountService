package models

import "github.com/dgrijalva/jwt-go"

// Claims for jwt token
type Claims struct {
	ID uint64
	jwt.StandardClaims
}
