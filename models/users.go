package models

import (
	"errors"

	"../../photofriends/hash"
	"../../photofriends/rand"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned when a resource
	// cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is
	// provided to a method like Delete
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password
	// is used when attempting to authenticate a user
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

const (
	userPwPepper  = "secret-random-string"
	hmacSecretKey = "secrey-hmac-key"
)

// User represent the user model stored in our database
// This is used for user accounts, storing both an email
// address and a password so users can log in and gain
// access to their individual private content
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"` // ignore in DB
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// UserDB is used to interact with the users database
//
// For all single user queries:
// 1 - user, nil 		- User found
// 2 - nil, ErrNotFound	- User not found
// 3 - nil, otherError  - Database error
type UserDB interface {

	// Methods for quering for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// used to close db connection
	Close() error

	// migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of methods used to mainpulate
// and work with the user model
type UserService interface {

	// Authenticate will verify the provided email and
	// password are correct, if correct, the user corresponding
	// to that email will be returned, if not the releated error
	// for the reason the method failed
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}

	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac:   hmac,
		UserDB: ug,
	}
	return &userService{
		UserDB: uv,
	}, nil
}

// ensure interface is matching
var _ UserService = &userService{}

// implementation of interface
type userService struct {
	UserDB
}

// Authenticate can be used to authenticate a user with the provided
// email address and password.
//	If the email address
// 		provided is invald, this will return nil, ErrNotFound
// 	If the password provided is invalid, this will return
// 		nil, ErrInvalidPassword
// 	If the email and password are both valid, this will return
//		user, nil
// 	Otherwise if another error is encountered this will return
// 		nil, error
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

/******************* VALIDATORS **************************/

// ensure interface is matching
var _ UserDB = &userGorm{}

type userValFunc func(*User) error

func runUsersValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}

	return nil
}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}

	if err := runUsersValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) Create(user *User) error {
	err := runUsersValFuncs(user,
		uv.bcryptPassword,
		uv.setRmemberIfUnset,
		uv.hmacRemember)

	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Update will hash a remember token if it is provided
func (uv *userValidator) Update(user *User) error {
	err := runUsersValFuncs(user, uv.bcryptPassword, uv.hmacRemember)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

//Delete will delete the user with the provided ID
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id

	err := runUsersValFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

// bcryptPassword will hash a users password with a
// predefined pepper(userPwPepper) and bcrypt if the
// password field is not the empty string
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}

	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes) // convert byteslice to string
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}

	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setRmemberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}

	token, err := rand.RememberToken()
	user.RememberHash = uv.hmac.Hash(token)
	if err != nil {
		return err
	}

	user.Remember = token
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrInvalidID
		}

		return nil
	})
}

/************************************************************/

// ensure interface is matching
var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	return &userGorm{
		db: db,
	}, nil
}

type userGorm struct {
	db *gorm.DB
}

// ByID will look up a user by the id provided
// Cases:
// 1 - user, nil 		- User found
// 2 - nil, ErrNotFound	- User not found
// 3 - nil, otherError  - Database error
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks up a user with the given email address
// and returns that user
//
// 1 - user, nil 		- User found
// 2 - nil, ErrNotFound	- User not found
// 3 - nil, otherError  - Database error
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with the given remember token
// and returns that user. This method expects the remember token
// to already be hashed. Err are the same as ByEmail
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create will create the provided user and backfill data
// like the ID, CreatedAt and UpdatedAt fields
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update will update the provided user with all of the
// data in the provided user object
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

//Delete  will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close closes the UserService database connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset drops the user table and rebuilds it
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}

	return ug.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table
func (ug *userGorm) AutoMigrate() error {
	err := ug.db.AutoMigrate(&User{}).Error
	return err
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
