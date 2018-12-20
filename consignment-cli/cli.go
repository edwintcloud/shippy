package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/edwintcloud/shippy/consignment-service/proto/consignment"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
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

	// Contact the server and print out its response
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	// parse file using our function above
	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("could not parse file: %v", err)
	}

	// create consignment using gGRPC server function
	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	// test by using the getconsignments method and printing results
	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
