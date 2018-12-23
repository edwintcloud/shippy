package main

import (
	pb "github.com/edwintcloud/shippy/user-service/proto/user"
	"github.com/jinzhu/gorm"
)

// Repository is our interface handler methods will conform to
type Repository interface {
	GetAll() ([]*pb.User, error)
	Get(id string) (*pb.User, error)
	Create(user *pb.User) error
	GetByEmail(email string) (*pb.User, error)
}

// UserRepository is our repo struct holding the db session reference
type UserRepository struct {
	db *gorm.DB
}

// Get retreives a user from the db by id
func (repo *UserRepository) Get(id string) (*pb.User, error) {
	var user *pb.User
	user.Id = id

	// find user in db and bind result to user struct
	if err := repo.db.Where("Id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	// return user if found and no error
	return user, nil
}

// GetAll retreives all users from the db
func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	var users []*pb.User

	// find users in db and bind results to user struct
	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}

	// return results if found and no error
	return users, nil
}

// GetByEmail finds a user in the db by email
func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	user := &pb.User{}

	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	// return user if found and password matches and no error
	return user, nil
}

// Create creates user in the db
func (repo *UserRepository) Create(user *pb.User) error {

	// create user in db
	if err := repo.db.Create(user).Error; err != nil {
		return err
	}

	// return no error
	return nil
}
