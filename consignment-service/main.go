// consignment-service/main.go
package main

import (
	"context"
	"fmt"

	// Import the generated protobuf code
	pb "github.com/edwintcloud/shippy/consignment-service/proto/consignment"

	micro "github.com/micro/go-micro"
)

// IRepository is our consignment interface
type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository is a dummy repository simulating the use of a datastore
// of some kind. This will be replaced with a real implementation later.
type Repository struct {
	consignments []*pb.Consignment
}

// Create saves our consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

// GetAll returns all consignments
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition.
type service struct {
	repo IRepository
}

// CreateConsignment is a create method handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

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
	repo := &Repository{}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// name must match package name in proto file
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	// Init will parse command line flags
	srv.Init()

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
