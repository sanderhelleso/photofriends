package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource
	// cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is
	// provided to a method like Delete
	ErrInvalidID = errors.New("models: ID provided was invalid")
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
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks up a user with the given email address
// and returns that user
//
// 1 - user, nil 		- User found
// 2 - nil, ErrNotFound	- User not found
// 3 - nil, otherError  - Database error
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User 
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// first will query using the provided gorm.DB and it
// will get the first item returned and place it into
// dst. If nothing is found in the query, it will 
// return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}

	return err
}

// will create the provided user and backfill data
// like the ID, CreatedAt and UpdatedAt fields
func (us *UserService) Create(user *User) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil { return err }

	user.PasswordHash = string(hashedBytes) // convert byteslice to string
	user.Password = ""
	return us.db.Create(user).Error
}

// will update the provided user with all of the
// data in the provided user object
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}

	user := User{ Model: gorm.Model{ ID: id }}
	return us.db.Delete(&user).Error
}

// closes the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// drops the user table and rebuilds it
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err;
	}

	return us.AutoMigrate()
}

// will attempt to automatically migrate the users table
func (us *UserService) AutoMigrate() error {
	err := us.db.AutoMigrate(&User{}).Error
	return err
}

type User struct {
	gorm.Model
	Name 		 string
	Email	  	 string `gorm:"not null;unique_index"`
	Password 	 string `gorm:"-"` // ignore in DB
	PasswordHash string `gorm:"not null"`
}