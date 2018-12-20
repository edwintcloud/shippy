package main

import (
	pb "github.com/edwintcloud/shippy/vessel-service/proto/vessel"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	dbName           = "shippy"
	vesselCollection = "vessels"
)

// Repository is our vessel interface
type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
	Create(vessel *pb.Vessel) error
	Close()
}

// VesselRepository is a dummy repository simulating the use of a datastore
// of some kind. This will be replaced with a real implementation later.
type VesselRepository struct {
	session *mgo.Session
}

// FindAvailable checks a specification against a map of vessels.
// If capacity and max weight are below a vessels capacity and max
// weight then return that vessel.
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	var vessel *pb.Vessel

	// Here we define a more complex query than our consignment-service's
	// GetAll function. Here we're asking for a vessel who's max weight and
	// capacity are greater than and equal to the given capacity and weight.
	// We're also using the `One` function here as that's all we want.
	err := repo.collection().Find(bson.M{
		"capacity":  bson.M{"$gte": spec.Capacity},
		"maxweight": bson.M{"$gte": spec.MaxWeight},
	}).One(&vessel)
	if err != nil {
		return nil, err
	}

	// return result
	return vessel, nil
}

// Create method creates a vessel specification
func (repo *VesselRepository) Create(vessel *pb.Vessel) error {
	return repo.collection().Insert(vessel)
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
func (repo *VesselRepository) Close() {
	repo.session.Close()
}

// here we create the new session for each query
func (repo *VesselRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(vesselCollection)
}
