package tests

import (
	"go-auth-app/controllers"
	"go-auth-app/middleware"
	"go-auth-app/services"
	"go-auth-app/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminForbidden(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	userService := &services.UserService{
		UserRepo: mockRepo,
	}

	controller := controllers.UserController{
		UserService: userService,
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Apply the AdminOnly middleware manually before the handler
	r.GET("/admin/users", middleware.AdminOnly(), func(c *gin.Context) {
		controller.GetAllUsers(c)
	})

	req, _ := http.NewRequest("GET", "/admin/users", nil)
	req.Header.Set("Content-Type", "application/json")

	// Simulate a user with the "user" role in the context
	// Since middleware.AdminOnly gets the role from the context,
	// simulate it using a wrapper - or set in a fake middleware
	r.Use(func(c *gin.Context) {
		c.Set("role", "user")
	})

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}
