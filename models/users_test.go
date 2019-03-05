package models

import (
	"fmt"
	"time"
	"testing"
)

 func testingUserService() (*UserService, error) {
	const (
		host 	 = "localhost"
		port 	 = 5432
		user	 = "postgres"
		password = "postgres"
		dbname 	 = "photofriends_test"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	us, err := NewUserService(psqlInfo)

	if err != nil {
		return nil, err
	}

	us.db.LogMode(false)

	// clear the users table between tests
	us.DestructiveReset()
	return us, nil
 }

 func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}

	user := User {
		Name: "Michael Scott",
		Email: "michael@dundermifflin.com",
	}

	err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Errorf("Expected ID > 0. Recieved %d", user.ID)
	}

	var timeoutAllowed = time.Duration(5 * time.Second);
	if time.Since(user.CreatedAt) > timeoutAllowed {
		t.Errorf("Expected CreatedAt to be recent. Recieved %s", user.CreatedAt)
	}

	if time.Since(user.UpdatedAt) > timeoutAllowed {
		t.Errorf("Expected CreatedAt to be recent. Recieved %s", user.UpdatedAt)
	}
 }