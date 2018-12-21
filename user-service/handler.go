package main

import (
	"context"
	"errors"
	"log"

	pb "github.com/edwintcloud/shippy/user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	repo         Repository
	tokenService Authable
}

// Get is our handler for retreiving a user from the repo
func (h *handler) Get(ctx context.Context, req *pb.User, res *pb.Response) error {

	// Get user from repo
	user, err := h.repo.Get(req.Id)
	if err != nil {
		return err
	}

	// set res and return no error
	res.User = user
	return nil
}

// GetAll is our handler for retreiving all users from the repo
func (h *handler) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {

	// Get users from repo
	users, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	// set res and return no error
	res.Users = users
	return nil
}

// Auth is our handler for authenticating a user from the repo
func (h *handler) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {

	// status message
	log.Println("Logging in user ", req.Email)

	// get user from repo by email
	user, err := h.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}

	// Compare given password with hashed password in db
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	// set password to "" so it isn't exposed in jwt
	user.Password = ""

	// generate jwt
	token, err := h.tokenService.Encode(user)
	if err != nil {
		return err
	}

	// set res and return no error
	res.Token = token
	return nil
}

// Create is our handler for creating a user in the repo
func (h *handler) Create(ctx context.Context, req *pb.User, res *pb.Response) error {

	// generate hashed password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return err
	}

	// set password to hashed password
	req.Password = string(hashedPass)

	// create user in repo
	if err := h.repo.Create(req); err != nil {
		return err
	}

	// set password to "" so it isn't exposed in the jwt
	req.Password = ""

	// generate jwt token
	token, err := h.tokenService.Encode(req)
	if err != nil {
		return err
	}

	// set res and return no error
	res.User = req
	res.Token = &pb.Token{Token: token}
	return nil
}

// ValidateToken is our handler for validating jwt token
func (h *handler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	// Decode token
	claims, err := h.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}

	// print claims
	log.Println(claims)

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	// set res and return no error
	res.Valid = true
	return nil
}
