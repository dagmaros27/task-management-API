package infrastructure_test

import (
	"net/http"
	"net/http/httptest"
	"task_managment_api/domain"
	"task_managment_api/infrastructure"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateUserToken(user domain.User) (string, domain.CustomError) {
	args := m.Called(user)
	return args.Get(0).(string), args.Get(1).(domain.CustomError)
}

func (m *MockJWTService) ValidateToken(tokenString string) (jwt.MapClaims, domain.CustomError) {
	args := m.Called(tokenString)
	return args.Get(0).(jwt.MapClaims), args.Get(1).(domain.CustomError)
}

type MiddlewareTestSuite struct {
	suite.Suite
	mockService *MockJWTService
	user        domain.User
	token       string
	authService infrastructure.AuthMiddlewareService
}

func (suite *MiddlewareTestSuite) SetupTest() {
	// Initialize JWT service and a mock user
	suite.mockService = new(MockJWTService)
	suite.user = domain.User{
		ID:       "user-id-123",
		Username: "testuser",
		Role:     "admin", // Set the role to "admin" for testing AdminMiddleware
	}
	suite.authService = infrastructure.NewAuthService(suite.mockService)

	// Stub the token generation and validation methods
	suite.mockService.On("GenerateUserToken", suite.user).Return("mocked-token", domain.CustomError{})
	suite.mockService.On("ValidateToken", "mocked-token").Return(jwt.MapClaims{
		"userId":   suite.user.ID,
		"username": suite.user.Username,
		"role":     suite.user.Role,
	}, domain.CustomError{})

	// Generate a valid JWT token for the user
	token, err := suite.mockService.GenerateUserToken(suite.user)
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

	middleware := suite.authService.AuthMiddleware()
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

	middleware := suite.authService.AuthMiddleware()
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

	middleware := suite.authService.AuthMiddleware()
	middleware(c)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.JSONEq(`{"message": "Authorization format must be Bearer {token}"}`, w.Body.String())
}

// TestAdminMiddlewareSuccess tests successful admin authorization
func (suite *MiddlewareTestSuite) TestAdminMiddlewareSuccess() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "admin") // Simulate setting the role from AuthMiddleware

	middleware := suite.authService.AdminMiddleware()
	middleware(c)

	suite.Equal(http.StatusOK, w.Code)
}

// TestAdminMiddlewareForbidden tests forbidden access for non-admins
func (suite *MiddlewareTestSuite) TestAdminMiddlewareForbidden() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "user") // Simulate setting a non-admin role

	middleware := suite.authService.AdminMiddleware()
	middleware(c)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.JSONEq(`{"message": "Admins only"}`, w.Body.String())
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
