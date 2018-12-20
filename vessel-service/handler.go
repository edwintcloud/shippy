package main

import (
	"context"

	pb "github.com/edwintcloud/shippy/vessel-service/proto/vessel"
	"github.com/globalsign/mgo"
)

// handler should implement all of the methods to satisfy the service
// we defined in our protobuf definition.
type handler struct {
	session *mgo.Session
}

// GetRepo clones db session to set repo
func (h *handler) GetRepo() Repository {
	return &VesselRepository{h.session.Clone()}
}

// FindAvailable here is handled by grpc service
func (h *handler) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {

	// getrepo and defer session to close after method completes
	repo := h.GetRepo()
	defer repo.Close()

	// Find the next available vessel
	vessel, err := repo.FindAvailable(req)
	if err != nil {
		return err
	}

	// Set the vessel as part of the response message type
	res.Vessel = vessel
	return nil
}

// Create method creates new vessel
func (h *handler) Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error {

	// getrepo and defer session to close after method completes
	repo := h.GetRepo()
	defer repo.Close()

	// create vessel
	if err := repo.Create(req); err != nil {
		return err
	}

	// Set the vessel to req data and created to true
	res.Vessel = req
	res.Created = true
	return nil
}
