// consignment-service/main.go
package main

import (
	"log"
	"net"

	// Import the generated protobuf code
	pb "shippy/consignment-service/proto/consignment"

	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Set port for gRPC server below
const port = ":50051"

// IRepository is our consignment interface
type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
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

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition.
type service struct {
	repo IRepository
}

// CreateConsignment is a create method handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {

	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Return match the `Response` message we created in the protobuf definition.
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}

	// Setup gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our implementation
	// into the auto-generated interface code for our protobuf definition
	pb.RegisterShippingServiceServer(s, &service{repo})

	// Register reflection service on gRPC server
	reflection.Register(s)
	log.Printf("Starting gRPC server on port%s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {

	}
}
