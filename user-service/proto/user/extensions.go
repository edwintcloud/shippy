package go_micro_srv_user

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

// BeforeCreate is a pre-save hook to create unique uuid
// for id field when creating a user.
func (model *User) BeforeCreate(scope *gorm.Scope) error {
	uuid, _ := uuid.NewV4()
	return scope.SetColumn("Id", uuid.String())
}
