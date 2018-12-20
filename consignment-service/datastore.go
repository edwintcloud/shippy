package main

import "github.com/globalsign/mgo"

// CreateSession creates the main session to our mongodb instance
func CreateSession(host string) (*mgo.Session, error) {

	// Connect to db
	session, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}

	// Monotonic consistence mode guarantees sequential reads and
	// writes while still managing to distribute some of the reading
	// load to secondary clusters. Faster than Strong mode and safer
	// than eventual mode which can't guarantee data will be in order.
	session.SetMode(mgo.Monotonic, true)

	// return session and nil error
	return session, nil
}
