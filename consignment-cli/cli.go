package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	pb "github.com/edwintcloud/shippy/consignment-service/proto/consignment"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/metadata"
)

const defaultFilename = "consignment.json"

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment

	// Read file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal json result to consignment struct
	json.Unmarshal(data, &consignment)

	// return result
	return consignment, err
}

func main() {

	cmd.Init()

	// Create new greeter client
	client := pb.NewShippingService("go.micro.srv.consignment", microclient.DefaultClient)

	// get command line args
	file := defaultFilename
	var token string
	fmt.Print("Enter token: ")
	fmt.Scanln(&token)
	fmt.Println(token)

	// parse file using our function above
	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("could not parse file: %v", err)
	}

	// create context which contains our given token
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	// create consignment using gGRPC server function
	r, err := client.CreateConsignment(ctx, consignment)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	// test by using the getconsignments method and printing results
	getAll, err := client.GetConsignments(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
