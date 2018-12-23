package main

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/edwintcloud/shippy/user-service/proto/user"
)

// Our jwt secret, this should be in an env file when deployed
var key = []byte("supersecretkey")

// CustomClaims is our struct to hold our jwt claims
type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

// Authable is the jwt token interface
type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

// TokenService is the struct that holds our repo
type TokenService struct {
	repo Repository
}

// Decode decodes a jwt token
func (s *TokenService) Decode(token string) (*CustomClaims, error) {

	// status text
	fmt.Println("Decoding token")

	// parse the token
	// ParseWithClaims takes the token, custom claims struct referfence,
	// and a function that returns the secret key as []byte
	tokenType, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// Validate token and return claims
	if claims, ok := tokenType.Claims.(*CustomClaims); ok && tokenType.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// Encode encodes data to a jwt token
func (s *TokenService) Encode(user *pb.User) (string, error) {

	// get expire time
	expires := time.Now().Add(time.Hour * 72).Unix()

	// Create claims
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expires,
			Issuer:    "go.micro.srv.user",
		},
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token and return
	return token.SignedString(key)
}
