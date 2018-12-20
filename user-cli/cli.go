package main

import (
	"context"
	"log"
	"os"

	pb "github.com/edwintcloud/shippy/user-service/proto/user"
	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
)

func main() {

	cmd.Init()

	// create new greeter client
	client := pb.NewUserService("go.micro.srv.user", microclient.DefaultClient)

	// Define our flags for command line args
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:  "name",
				Usage: "Your full name",
			},
			cli.StringFlag{
				Name:  "email",
				Usage: "Your email",
			},
			cli.StringFlag{
				Name:  "password",
				Usage: "Your password",
			},
			cli.StringFlag{
				Name:  "company",
				Usage: "Your company",
			},
		),
	)

	// Start as a service
	service.Init(
		micro.Action(func(c *cli.Context) {

			// Get our command line arg values
			name := c.String("name")
			email := c.String("email")
			password := c.String("password")
			company := c.String("company")

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

			// Quit client
			os.Exit(0)
		}),
	)

	// Run cli service
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
