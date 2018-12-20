package main

import (
	pb "github.com/edwintcloud/shippy/consignment-service/proto/consignment"
	"github.com/globalsign/mgo"
)

const (
	dbName                = "shippy"
	consignmentCollection = "consignments"
)

// Repository is our consignment interface
type Repository interface {
	Create(*pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
	Close()
}

// ConsignmentRepository is a dummy repository simulating the use of a datastore
// of some kind. This will be replaced with a real implementation later.
type ConsignmentRepository struct {
	session *mgo.Session
}

// Create saves our consignment
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) error {
	return repo.collection().Insert(consignment)
}

// GetAll returns all consignments
func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
	var consignments []*pb.Consignment

	// Get all documents from db and bind to consignments var
	err := repo.collection().Find(nil).All(&consignments)

	// return results
	return consignments, err
}

// Close closes the database session after each query has ran.
// Mgo creates a 'master' session on start-up, it's then good practice
// to copy a new session for each request that's made. This means that
// each request has its own database session. This is safer and more efficient,
// as under the hood each session has its own database socket and error handling.
// Using one main database socket means requests having to wait for that session.
// I.e this approach avoids locking and allows for requests to be processed concurrently. Nice!
// But... it does mean we need to ensure each session is closed on completion. Otherwise
// you'll likely build up loads of dud connections and hit a connection limit. Not nice!
func (repo *ConsignmentRepository) Close() {
	repo.session.Close()
}

// here we create the new session for each query
func (repo *ConsignmentRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(consignmentCollection)
}
