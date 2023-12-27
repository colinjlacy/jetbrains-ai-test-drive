package services

import (
	"github.com/colinjlacy/jetbrains-ai-test-drive/models"
	"testing"
)

func mockUsers() {
	users = map[string]*models.User{
		"1": &models.User{Id: "1", Name: "Mario", Age: 38},
		"2": &models.User{Id: "2", Name: "Luigi", Age: 35},
		"3": &models.User{Id: "3", Name: "Peach", Age: 37},
		"4": &models.User{Id: "4", Name: "Toad", Age: 73},
	}
}

func TestDeleteUserById(t *testing.T) {

	t.Run("returns error when user not found", func(t *testing.T) {
		mockUsers()
		err := DeleteUserById("nonexistent")
		if err != ErrorUserNotFound {
			t.Errorf("Expected error %q, but got %q", ErrorUserNotFound, err)
		}
	})

	t.Run("deletes user by id", func(t *testing.T) {
		mockUsers()
		err := DeleteUserById("1")
		if err != nil {
			t.Errorf("Expected no error, but got %q", err)
		}

		if _, exists := users["1"]; exists {
			t.Errorf("User was not deleted")
		}
	})
}

// Testing the UpsertUser function
func TestUpsertUser(t *testing.T) {

	t.Run("returns error when user fields are nil", func(t *testing.T) {
		mockUsers()
		user := models.User{Id: "", Name: "", Age: 0}
		err := UpsertUser(&user)
		if err != ErrorUserFieldNil {
			t.Errorf("Expected error %q, but got %q", ErrorUserFieldNil, err)
		}
	})

	t.Run("returns error when new user's name is not unique", func(t *testing.T) {
		mockUsers()
		user := models.User{Id: "23", Name: "Mario", Age: 2}
		err := UpsertUser(&user)
		if err != ErrorUserNameExists {
			t.Errorf("Expected error %q, but got %q", ErrorUserNameExists, err)
		}
	})

	t.Run("updates existing user in map", func(t *testing.T) {
		mockUsers()
		user := models.User{Id: "1", Name: "Steve", Age: 42}
		err := UpsertUser(&user)
		if err != nil {
			t.Errorf("Expected no error, but got %q", err)
		}

		if users[user.Id].Name != "Steve" {
			t.Errorf("User was not updated in map")
		}
	})

	t.Run("adds new user to map", func(t *testing.T) {
		mockUsers()
		user := models.User{Id: "5", Name: "Robin", Age: 30}
		err := UpsertUser(&user)
		if err != nil {
			t.Errorf("Expected no error, but got %q", err)
		}

		if users[user.Id] != &user {
			t.Errorf("User was not added to map")
		}
	})
}

func TestCreateUser(t *testing.T) {

	t.Run("returns error when user id exists", func(t *testing.T) {
		mockUsers()
		user := models.User{Id: "1", Name: "Bowser", Age: 40}
		err := CreateUser(&user)
		if err != ErrorUserExists {
			t.Errorf("Expected error %q, but got %q", ErrorUserExists, err)
		}
	})

	t.Run("returns error when user name exists", func(t *testing.T) {
		mockUsers()
		user := models.User{Id: "5", Name: "Mario", Age: 40}
		err := CreateUser(&user)
		if err != ErrorUserNameExists {
			t.Errorf("Expected error %q, but got %q", ErrorUserNameExists, err)
		}
	})

	t.Run("adds new user to map", func(t *testing.T) {
		mockUsers()
		user := models.User{Id: "5", Name: "Bowser", Age: 40}
		err := CreateUser(&user)
		if err != nil {
			t.Errorf("Expected no error, but got %q", err)
		}

		if users[user.Id] != &user {
			t.Errorf("User was not added to map")
		}
	})
}

func TestGetUserById(t *testing.T) {

	t.Run("returns the correct user", func(t *testing.T) {
		mockUsers()

		user := GetUserById("1")

		if user != users["1"] {
			t.Errorf("GetUserById() failed, got %v, want %v", user, users["1"])
		}
	})

	t.Run("returns nil for non existent user", func(t *testing.T) {
		mockUsers()

		user := GetUserById("0")

		if user != nil {
			t.Errorf("GetUserById() failed, expected a nil, got: %v", user)
		}
	})
}

func TestGetUsers(t *testing.T) {
	mockUsers()
	userList := GetUsers()
	if len(userList) != 4 {
		t.Errorf("GetUsers() failed, expected length of 4, got %d", len(userList))
	}
}
