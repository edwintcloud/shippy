package main

import (
	"fmt"
	"log"

	pb "github.com/edwintcloud/shippy/user-service/proto/user"
	"github.com/micro/go-micro"
)

func main() {

	// Create db connection and defer db to close when method ends
	db, err := CreateConnection()
	defer db.Close()
	if err != nil {
		log.Fatalf("could not connect to DB: %v", err)
	}

	// AutoMigrate database to match pb User struct
	db.AutoMigrate(&pb.User{})

	// set UserRepository to use db conection
	repo := &UserRepository{db}

	// set TokenService to use repo
	tokenService := &TokenService{repo}

	// Create new micro service
	srv := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// Initialize service with command line args
	srv.Init()

	// Register service handler
	pb.RegisterUserServiceHandler(srv.Server(), &handler{repo, tokenService})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
