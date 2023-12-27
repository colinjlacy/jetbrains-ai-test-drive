package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	t.Run("GET root path", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		router := gin.Default()

		router.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hello, world",
			})
		})

		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatalf("Couldn't create request: %s\n", err)
		}

		respRec := httptest.NewRecorder()
		router.ServeHTTP(respRec, request)

		assert.Equal(t, http.StatusOK, respRec.Code)
		assert.Equal(t, `{"message":"hello, world"}`, respRec.Body.String())
	})
}

func TestRegisterHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	registerHealthEndpoint(r)

	req, _ := http.NewRequest("GET", "/health", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.Code)
	}

	expectedResponse := `{"status":"ok"}`
	body := resp.Body.String()
	if body != expectedResponse {
		t.Errorf("Expected body '%s', but got '%s'", expectedResponse, body)
	}
}
