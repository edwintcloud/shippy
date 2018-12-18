package main

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/edwintcloud/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
)

// Repository is the interface we will use to define our methods
type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

// VesselRepository is a struct that will hold vessels
type VesselRepository struct {
	vessels []*pb.Vessel
}

// FindAvailable checks a specification against a map of vessels.
// If capacity and max weight are below a vessels capacity and max
// weight then return that vessel.
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("no vessel found by that spec")
}

// grpc service handler
type service struct {
	repo Repository
}

// FindAvailable here is handled by grpc service
func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {

	// Find the next available vessel
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}

	// Set the vessel as part of the response message type
	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessels}

	// create micro service
	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)

	// Initialize service with command line args
	srv.Init()

	// Register handler and run service
	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

}
