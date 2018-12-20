package main

import (
	"context"

	pb "github.com/edwintcloud/shippy/user-service/proto/user"
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

	// authenticate user from repo
	_, err := h.repo.GetByEmailAndPassword(req)
	if err != nil {
		return err
	}

	// set res and return no error
	res.Token = "testingabcc"
	return nil
}

// Create is our handler for creating a user in the repo
func (h *handler) Create(ctx context.Context, req *pb.User, res *pb.Response) error {

	// create user in repo
	if err := h.repo.Create(req); err != nil {
		return err
	}

	// set res and return no error
	res.User = req
	return nil
}

// ValidateToken is our handler for validating jwt token
func (h *handler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	// TODO: Implement
	return nil
}
