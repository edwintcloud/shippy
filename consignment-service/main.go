// consignment-service/main.go
package main

import (
	"context"
	"fmt"

	// Import the generated protobuf code
	pb "github.com/edwintcloud/shippy/consignment-service/proto/consignment"
	vesselPb "github.com/edwintcloud/shippy/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
)

// Repository is our consignment interface
type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// ConsignmentRepository is a dummy repository simulating the use of a datastore
// of some kind. This will be replaced with a real implementation later.
type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

// Create saves our consignment
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

// GetAll returns all consignments
func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition.
type service struct {
	repo         Repository
	vesselClient vesselPb.VesselService
}

// CreateConsignment is a create method handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Call client instance of our vessel service with our consignment weight,
	// and the amount of container as the capacity value
	vesselRes, err := s.vesselClient.FindAvailable(context.Background(), &vesselPb.Specification{
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
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// Return match the `Response` message we created in the protobuf definition.
	res.Created = true
	res.Consignment = consignment
	return nil
}

// GetConsignments is the getrequest method handled by the gRPC server.
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &ConsignmentRepository{}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// name must match package name in proto file
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	// Create vessel service client
	vesselClient := vesselPb.NewVesselService("go.micro.srv.vessel", srv.Client())

	// Init will parse command line flags
	srv.Init()

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
