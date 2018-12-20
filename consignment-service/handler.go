package main

import (
	"context"
	"fmt"

	pb "github.com/edwintcloud/shippy/consignment-service/proto/consignment"
	vesselPb "github.com/edwintcloud/shippy/vessel-service/proto/vessel"
	"github.com/globalsign/mgo"
)

// handler should implement all of the methods to satisfy the service
// we defined in our protobuf definition.
type handler struct {
	session       *mgo.Session
	vesselService vesselPb.VesselService
}

// GetRepo clones db session to set repo
func (h *handler) GetRepo() Repository {
	return &ConsignmentRepository{h.session.Clone()}
}

// CreateConsignment is a create method handled by the gRPC server.
func (h *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// get repo and defer session to close after method completion
	repo := h.GetRepo()
	defer repo.Close()

	// Call client instance of our vessel service with our consignment weight,
	// and the amount of container as the capacity value
	vesselRes, err := h.vesselService.FindAvailable(context.Background(), &vesselPb.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	fmt.Printf("Found vessel: %s \n", vesselRes.Vessel.Name)
	if err != nil {
		return err
	}

	// Set the VesselId as the vessel response we got from vessel service
	req.VesselId = vesselRes.Vessel.Id

	// Save our consignment
	err = repo.Create(req)
	if err != nil {
		return err
	}

	// Return match the `Response` message we created in the protobuf definition.
	res.Created = true
	res.Consignment = req
	return nil
}

// GetConsignments is the getrequest method handled by the gRPC server.
func (h *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {

	// get repo and defer session to close after method completes
	repo := h.GetRepo()
	defer repo.Close()

	// Get all consignments
	consignments, err := repo.GetAll()
	if err != nil {
		return err
	}

	// return results
	res.Consignments = consignments
	return nil
}
