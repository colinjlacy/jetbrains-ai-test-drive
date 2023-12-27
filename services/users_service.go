package services

import (
	"errors"
	"github.com/colinjlacy/jetbrains-ai-test-drive/models"
)

var users = map[string]*models.User{
	"1": &models.User{Id: "1", Name: "Mario", Age: 38},
	"2": &models.User{Id: "2", Name: "Luigi", Age: 35},
	"3": &models.User{Id: "3", Name: "Peach", Age: 37},
	"4": &models.User{Id: "4", Name: "Toad", Age: 73},
}

var ErrorUserExists = errors.New("user already exists")
var ErrorUserNameExists = errors.New("user with this name already exists")
var ErrorUserFieldNil = errors.New("user fields cannot be nil")
var ErrorUserNotFound = errors.New("user not found")

func GetUserById(id string) *models.User {
	return users[id]
}

func GetUsers() []models.User {
	userList := make([]models.User, 0, len(users))
	for _, user := range users {
		userList = append(userList, *user)
	}
	return userList
}

func CreateUser(user *models.User) error {
	_, idExists := users[user.Id]
	for _, u := range users {
		println("u.Name", u.Name, "user.Name", user.Name)
		if u.Name == user.Name {
			return ErrorUserNameExists
		}
	}

	if idExists {
		return ErrorUserExists
	}

	users[user.Id] = user
	return nil
}

func UpsertUser(user *models.User) error {
	if user == nil || user.Id == "" || user.Name == "" {
		return ErrorUserFieldNil
	}

	// Check for uniqueness of Name and Id
	if _, exists := users[user.Id]; !exists {
		for id, u := range users {
			if u.Name == user.Name && id != user.Id {
				return ErrorUserNameExists
			}
		}
	}

	// Upsert operation
	users[user.Id] = user
	return nil
}

func DeleteUserById(id string) error {
	_, ok := users[id]
	if !ok {
		return ErrorUserNotFound
	}
	delete(users, id)
	return nil
}
