package endpoints

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/colinjlacy/jetbrains-ai-test-drive/models"
	"github.com/colinjlacy/jetbrains-ai-test-drive/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	// imported errors
	errorUserExists     error = services.ErrorUserExists
	errorUserNameExists error = services.ErrorUserNameExists
	errorUserFieldNil   error = services.ErrorUserFieldNil
	errorUserNotFound   error = services.ErrorUserNotFound

	// factory functions for mocks
	createUserMock func(err error) func(user *models.User) error = func(err error) func(user *models.User) error {
		return func(user *models.User) error {
			return err
		}
	}

	upsertUserMock func(err error) func(user *models.User) error = func(err error) func(user *models.User) error {
		return func(user *models.User) error {
			return err
		}
	}

	deleteUserByIdMock func(err error) func(id string) error = func(err error) func(id string) error {
		return func(id string) error {
			return err
		}
	}
	getUsersMock func(success bool) func() []models.User = func(success bool) func() []models.User {
		if success {
			return func() []models.User {
				return []models.User{
					{Id: "1", Name: "Dummy User 1", Age: 30},
					{Id: "2", Name: "Dummy User 2", Age: 60},
				}
			}
		} else {
			return func() []models.User {
				return nil
			}
		}
	}

	getUserByIdMock func(success bool) func(id string) *models.User = func(success bool) func(id string) *models.User {
		if success {
			return func(id string) *models.User {
				return &models.User{Id: id, Name: "Dummy User " + id, Age: 45}
			}
		} else {
			return func(id string) *models.User {
				return nil
			}
		}
	}
)

func TestGetUsersHandler_WithMultipleEntries(t *testing.T) {
	getUsers = getUsersMock(true)

	// Prepare HTTP route with gin context.
	router := gin.Default()
	router.GET("/users", getUsersHandler)
	httpResponse := performGetRequest(router, "/users")

	// Check for correct http response.
	if httpResponse.Code != http.StatusOK {
		t.Fatalf("Expected the http response code to be %v, but got %v", http.StatusOK, httpResponse.Code)
	}
}

func TestGetUsersHandler_WithNoEntries(t *testing.T) {
	getUsers = getUsersMock(false)

	// Prepare HTTP route with gin context.
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users", getUsersHandler)
	httpResponse := performGetRequest(router, "/users")

	// Check for correct http response.
	if httpResponse.Code != http.StatusOK {
		t.Fatalf("Expected the http response code to be %v, but got %v", http.StatusOK, httpResponse.Code)
	}
}

func TestGetUserByIdHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/users/:id", getUserByIdHandler)

	// Arrange: Setting up the positive test
	getUserById = func(id string) *models.User {
		return &models.User{Id: "1", Name: "John Doe", Age: 30}
	}
	request, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	responseRecorder := httptest.NewRecorder()

	// Act: Make the request
	r.ServeHTTP(responseRecorder, request)

	// Assert: Check if the status code is 200 and payload is as expected
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "John Doe")

	// Arrange: Setting up the test for a case when user doesn't exist
	getUserById = func(id string) *models.User {
		return nil
	}
	request2, _ := http.NewRequest(http.MethodGet, "/users/2", nil)
	responseRecorder2 := httptest.NewRecorder()

	// Act: Make the request
	r.ServeHTTP(responseRecorder2, request2)

	// Assert: Check if the status code is 404
	assert.Equal(t, http.StatusNotFound, responseRecorder2.Code)
}

func TestCreateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	t.Run("createUser returns error", func(t *testing.T) {
		createUser = func(*models.User) error { return errors.New("mock error") }

		router.POST("/user", createUserHandler)

		req := httptest.NewRequest(http.MethodPost, "/user", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, resp.Code)
		}
	})

	t.Run("createUser does not return error", func(t *testing.T) {
		router = gin.New()
		createUser = func(*models.User) error { return nil }

		router.POST("/user", createUserHandler)

		req := httptest.NewRequest(http.MethodPost, "/user", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Fatalf("Expected status code %d, got %d", http.StatusCreated, resp.Code)
		}
	})

	t.Run("createUser returns error when failing to bind", func(t *testing.T) {
		router = gin.New()
		createUser = func(*models.User) error { return nil }

		router.POST("/user", createUserHandler)

		req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader("test"))
		req.Header.Set("Content-Type", "application/xml")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, resp.Code)
		}
	})
}

type upsertUserHandlerTestSuite struct {
	suite.Suite
	ctx  *gin.Context
	w    *httptest.ResponseRecorder
	user models.User
}

func (suite *upsertUserHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.w = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.w)

	suite.user = models.User{Id: "1", Name: "Dummy User 1", Age: 30}

	userJson, _ := json.Marshal(suite.user)
	req, _ := http.NewRequest(http.MethodPut, "/user/1", bytes.NewBuffer(userJson))
	suite.ctx.Request = req
	suite.ctx.Params = gin.Params{gin.Param{Key: "user-id", Value: "1"}}
}

func TestUpsertUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(upsertUserHandlerTestSuite))
}

func (suite *upsertUserHandlerTestSuite) TestupsertUserHandlerSuccess() {
	mockService := new(MockService)
	mockService.On("upsertUserMock", mock.AnythingOfType("*models.User")).Return(nil)

	oldFunc := upsertUser
	upsertUser = mockService.upsertUserMock

	upsertUserHandler(suite.ctx)

	mockService.AssertExpectations(suite.T())
	suite.Require().Equal(http.StatusOK, suite.w.Code)

	upsertUser = oldFunc
}

func (suite *upsertUserHandlerTestSuite) TestupsertUserHandlerError() {
	mockService := new(MockService)
	mockService.On("upsertUserMock", mock.AnythingOfType("*models.User")).Return(services.ErrorUserNotFound)

	oldFunc := upsertUser
	upsertUser = mockService.upsertUserMock

	upsertUserHandler(suite.ctx)

	mockService.AssertExpectations(suite.T())
	suite.Require().Equal(http.StatusBadRequest, suite.w.Code)

	upsertUser = oldFunc
}

func (suite *upsertUserHandlerTestSuite) TestupsertUserHandlerBindingError() {
	req, _ := http.NewRequest(http.MethodPut, "/user/1", strings.NewReader("test"))
	suite.ctx.Request = req

	upsertUserHandler(suite.ctx)

	suite.Require().Equal(http.StatusBadRequest, suite.w.Code)
}

type MockService struct {
	mock.Mock
}

func (m *MockService) upsertUserMock(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Utility function to perform a HTTP GET request.
func performGetRequest(r http.Handler, route string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", route, nil)
	responseRecorder := httptest.NewRecorder()
	r.ServeHTTP(responseRecorder, req)
	return responseRecorder
}

// Now we need test deleteUserHandler function
func TestDeleteUserHandler(t *testing.T) {
	// testing case when user was deleted
	deleteUserById = deleteUserByIdMock(nil)
	router := gin.Default()
	router.DELETE("/user/:user-id", deleteUserHandler)

	req, _ := http.NewRequest(http.MethodDelete, "/user/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if status := resp.Result().StatusCode; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// testing case when user not found
	deleteUserById = deleteUserByIdMock(errorUserNotFound)
	req, _ = http.NewRequest(http.MethodDelete, "/user/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if status := resp.Result().StatusCode; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// testing case when user not found
	deleteUserById = deleteUserByIdMock(errors.New("test error"))
	req, _ = http.NewRequest(http.MethodDelete, "/user/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if status := resp.Result().StatusCode; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestRegisterUserEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := gin.New()
	RegisterUserEndpoints(e)
	if len(e.Routes()) > 0 == false {
		t.Errorf("expected gin engine to have registered routes, but found none")
	}
}
