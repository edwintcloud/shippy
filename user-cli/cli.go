package main

import (
	"context"
	"log"
	"os"

	pb "github.com/edwintcloud/shippy/user-service/proto/user"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
)

func main() {

	cmd.Init()

	// create new greeter client
	client := pb.NewUserService("go.micro.srv.user", microclient.DefaultClient)

	// Test user
	name := "Edwin Cloud"
	email := "ecloud412@gmail.com"
	password := "password"
	company := "ETC"

	// call user service to create user
	r, err := client.Create(context.TODO(), &pb.User{
		Name:     name,
		Email:    email,
		Password: password,
		Company:  company,
	})
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}

	// print result
	log.Printf("Created user with id %s", r.User.Id)

	// call user service to get all users
	getAll, err := client.GetAll(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("unable to get all users: %v", err)
	}

	// print results
	for _, v := range getAll.Users {
		log.Println(v)
	}

	// Test auth
	authResponse, err := client.Auth(context.TODO(), &pb.User{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatalf("could not authenticate user: %s error: %v\n", email, err)
	}

	// print auth result
	log.Printf("Your access token is: %s\n", authResponse.Token)

	// Quit client
	os.Exit(0)
}
