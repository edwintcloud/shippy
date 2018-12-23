// consignment-service/main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	// Import the generated protobuf code
	pb "github.com/edwintcloud/shippy/consignment-service/proto/consignment"
	userPb "github.com/edwintcloud/shippy/user-service/proto/user"
	vesselPb "github.com/edwintcloud/shippy/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
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
		// Auth middleware
		micro.WrapHandler(AuthWrapper),
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

// AuthWrapper is a high-order function which takes a HandlerFunc
// and returns a function, which takes a context, request and response interface.
// The token is extracted from the context set in our consignment-cli, that
// token is then sent over to the user service to be validated.
// If valid, the call is passed along to the handler. If not,
// an error is returned.
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {

		// if disable auth true then return without decode
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		}

		// get metadata from context
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}
		token := meta["Token"]

		// print a status
		log.Println("Authenticating with token: ", token)

		// validate token using user service
		authClient := userPb.NewUserService("go.micro.srv.user", microclient.DefaultClient)
		_, err := authClient.ValidateToken(ctx, &userPb.Token{
			Token: token,
		})
		if err != nil {
			return err
		}

		// return result
		err = fn(ctx, req, resp)
		return err
	}
}
