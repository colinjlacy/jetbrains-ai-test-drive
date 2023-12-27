package endpoints

import (
	"github.com/colinjlacy/jetbrains-ai-test-drive/models"
	"github.com/colinjlacy/jetbrains-ai-test-drive/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	getUsers       func() []models.User          = services.GetUsers
	getUserById    func(id string) *models.User  = services.GetUserById
	createUser     func(user *models.User) error = services.CreateUser
	upsertUser     func(user *models.User) error = services.UpsertUser
	deleteUserById func(id string) error         = services.DeleteUserById
)

func getUsersHandler(c *gin.Context) {
	users := getUsers()
	c.JSON(http.StatusOK, users)
}

func getUserByIdHandler(c *gin.Context) {
	id := c.Param("user-id")
	user := getUserById(id)
	if user != nil {
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "User not found",
		})
	}
}

func createUserHandler(c *gin.Context) {
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := createUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func upsertUserHandler(context *gin.Context) {
	var user models.User
	if err := context.BindJSON(&user); err == nil {
		if userId, exists := context.Params.Get("user-id"); exists {
			user.Id = userId
			if err := upsertUser(&user); err == nil {
				context.JSON(http.StatusOK, gin.H{})
			} else {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
	}
}

func deleteUserHandler(c *gin.Context) {
	id := c.Param("user-id")
	err := deleteUserById(id)

	if err != nil {
		if err == services.ErrorUserNotFound {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusNoContent)
}
