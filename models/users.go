package models

import (
	"errors"
	"strings"

	"regexp"

	"../../photofriends/hash"
	"../../photofriends/rand"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned when a resource
	// cannot be found in the database
	ErrNotFound = errors.New("Resource not found")

	// ErrIDInvalid is returned when an invalid ID is
	// provided to a method like Delete
	ErrIDInvalid = errors.New("ID provided was invalid")

	// ErrPasswordIncorrect is returned when an invalid password
	// is used when attempting to authenticate a user
	ErrPasswordIncorrect = errors.New("Incorrect password provided")

	// ErrEmailRequired is returned when an email address
	// is not provided when creating a user
	ErrEmailRequired = errors.New("Email address is required")

	// ErrEmailInvalid is returned when an email address provided
	// does not match any of our requirements
	ErrEmailInvalid = errors.New("Email address is not valid")

	// ErrEmailTaken is returned when an update or create
	// is attempted with an email address that is already in use
	ErrEmailTaken = errors.New("Email address is already taken")

	// ErrPasswordTooShort is returned when an update or create is
	// attempted with a user passord that is less than 8 characters
	ErrPasswordTooShort = errors.New("Password must be atleast 8 characters long")

	// ErrPasswordRequired is returned when a create is attempted
	// wihtout a user password provided
	ErrPasswordRequired = errors.New("Password is required")

	// ErrRememberRequired is returned when a create or update
	// is attempted wihtout a user remember token hash provided
	ErrRememberRequired = errors.New("Remember token is required")

	// ErrRememberTooShort is returned when a remember token is
	// not atleast 32 bytes of length
	ErrRememberTooShort = errors.New("Remember token must be atleast 32 bytes")
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

func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)

	return &userService{
		UserDB: uv,
	}
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
// 		nil, ErrPasswordIncorrect
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
			return nil, ErrPasswordIncorrect
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

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB: udb,
		hmac:   hmac,

		// emailRegex is used to match email addresses.
		// It is not perfect, but works well enough for now
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

// ByEmail will normalize the email address before
// calling ByEmail on the UserDB field
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}

	err := runUsersValFuncs(&user,
		uv.normalizeEmail,
		uv.emailFormat,
		uv.requireEmail)

	if err != nil {
		return nil, err
	}

	return uv.ByEmail(user.Email)
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
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRmemberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.emailFormat,
		uv.requireEmail,
		uv.emailFormat)

	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Update will hash a remember token if it is provided
func (uv *userValidator) Update(user *User) error {
	err := runUsersValFuncs(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.emailFormat,
		uv.requireEmail,
		uv.emailIsAvail)

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

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}

	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}

	if n < 32 {
		return ErrRememberTooShort
	}

	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}

	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrIDInvalid
		}

		return nil
	})
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}

	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}

	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)

	if err == ErrNotFound {
		// email address is not taken
		return nil
	}

	if err != nil {
		return err
	}

	// we found a user with this email address...
	// If the found user has the same ID as this user,
	// it is an update and this is the same user of email
	if user.ID != existing.ID {
		return ErrEmailTaken
	}

	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}

	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrEmailInvalid
	}

	return nil
}

/************************************************************/

// ensure interface is matching
var _ UserDB = &userGorm{}

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
