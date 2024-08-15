package infrastructure_test

import (
	"net/http"
	"net/http/httptest"
	"task_managment_api/domain"
	"task_managment_api/infrastructure"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
	suite.Suite
	jwtService infrastructure.JWTService
	user       domain.User
	token      string
}

func (suite *MiddlewareTestSuite) SetupTest() {
	// Initialize JWT service and a mock user
	suite.jwtService = infrastructure.NewJWTService()
	suite.user = domain.User{
		ID:       "user-id-123",
		Username: "testuser",
		Role:     "admin", // Set the role to "admin" for testing AdminMiddleware
	}

	// Generate a valid JWT token for the user
	token, err := suite.jwtService.GenerateUserToken(suite.user)
	if err.ErrCode != 0 {
		suite.T().Fatal("Failed to generate token:", err.ErrMessage)
	}
	suite.token = token
}

// TestAuthMiddlewareSuccess tests successful authorization
func (suite *MiddlewareTestSuite) TestAuthMiddlewareSuccess() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+suite.token)

	middleware := infrastructure.AuthMiddleware()
	middleware(c)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(suite.user.ID, c.MustGet("userId"))
	suite.Equal(suite.user.Username, c.MustGet("username"))
	suite.Equal(suite.user.Role, c.MustGet("role"))
}

// TestAuthMiddlewareMissingAuthorizationHeader tests missing authorization header
func (suite *MiddlewareTestSuite) TestAuthMiddlewareMissingAuthorizationHeader() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	middleware := infrastructure.AuthMiddleware()
	middleware(c)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.JSONEq(`{"message": "Authorization header required"}`, w.Body.String())
}

// TestAuthMiddlewareInvalidFormat tests invalid authorization format
func (suite *MiddlewareTestSuite) TestAuthMiddlewareInvalidFormat() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "InvalidTokenFormat")

	middleware := infrastructure.AuthMiddleware()
	middleware(c)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.JSONEq(`{"message": "Authorization format must be Bearer {token}"}`, w.Body.String())
}

// TestAdminMiddlewareSuccess tests successful admin authorization
func (suite *MiddlewareTestSuite) TestAdminMiddlewareSuccess() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "admin") // Simulate setting the role from AuthMiddleware

	middleware := infrastructure.AdminMiddleware()
	middleware(c)

	suite.Equal(http.StatusOK, w.Code)
}

// TestAdminMiddlewareForbidden tests forbidden access for non-admins
func (suite *MiddlewareTestSuite) TestAdminMiddlewareForbidden() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "user") // Simulate setting a non-admin role

	middleware := infrastructure.AdminMiddleware()
	middleware(c)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.JSONEq(`{"message": "Admins only"}`, w.Body.String())
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
