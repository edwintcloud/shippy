// consignment-service/main.go
package main

import (
	"fmt"
	"log"
	"os"

	// Import the generated protobuf code
	pb "github.com/edwintcloud/shippy/consignment-service/proto/consignment"
	vesselPb "github.com/edwintcloud/shippy/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
)

// default mongo database connection uri
const defaultHost = "localhost:27017"

func main() {

	// try to get mongodb connection uri from env variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = defaultHost
	}

	// Create database session
	session, err := CreateSession(host)

	// defer session close so database connection is closed on main func exit
	defer session.Close()

	if err != nil {
		log.Panicf("could not connect to datastore with host %s - %v", host, err)
	}

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
	pb.RegisterShippingServiceHandler(srv.Server(), &handler{session, vesselClient})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
