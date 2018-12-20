package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/edwintcloud/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
)

// defaultHost is our mongoDB default connection uri
const defaultHost = "localhost:27017"

// createDummyData creates dummy data to test our implementation
func (repo *VesselRepository) createDummyData() {

	// defer db session to close after method completion
	defer repo.Close()

	// create slice of vessels
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Kane's Salty Secret", MaxWeight: 200000, Capacity: 500},
	}

	// write vessels to db using handler Create method
	for _, v := range vessels {
		repo.Create(v)
	}
}

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

	// copy session as repo
	repo := &VesselRepository{session.Copy()}

	// test our service by using the dummy data func above
	repo.createDummyData()

	// create micro service
	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)

	// Initialize service with command line args
	srv.Init()

	// Register handler and run service
	pb.RegisterVesselServiceHandler(srv.Server(), &handler{session})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

}
