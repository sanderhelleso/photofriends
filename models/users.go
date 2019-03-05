package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource
	// cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")
)

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {	
		return nil, err
	}

	return &UserService {
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

// ByID will look up a user by the id provided
// Cases:
// 1 - user, nil 		- User found
// 2 - nil, ErrNotFound	- User not found
// 3 - nil, otherError  - Database error
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}	
}

// will create the provided user and backfill data
// like the ID, CreatedAt and UpdatedAt fields
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// will update the provided user with all of the
// data in the provided user object
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// closes the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// drops the user table and rebuilds it
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name string
	Email string `gorm:"not null;unique_index"`
}