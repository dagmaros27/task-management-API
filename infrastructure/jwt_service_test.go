package infrastructure_test

import (
	"net/http"
	bootstrap "task_managment_api"
	"task_managment_api/domain"
	"task_managment_api/infrastructure"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/suite"
)

type JWTServiceTestSuite struct {
	suite.Suite
	service infrastructure.JWTService
	secret  string
	user    domain.User
}

func (suite *JWTServiceTestSuite) SetupTest() {
	suite.service = infrastructure.NewJWTService()
	suite.secret = "testsecret" // Mocked secret key for testing
	suite.user = domain.User{
		ID:       "user-id-123",
		Username: "testuser",
		Role:     "user",
	}

	// Mocking bootstrap.NewEnv().AccessTokenSecret
	bootstrapEnv := bootstrap.NewEnv()
	bootstrapEnv.AccessTokenSecret = suite.secret
}

// TestGenerateUserTokenSuccess tests successful token generation
func (suite *JWTServiceTestSuite) TestGenerateUserTokenSuccess() {
	token, err := suite.service.GenerateUserToken(suite.user)

	suite.Empty(err.ErrCode)
	suite.NotEmpty(token)
}

// TestValidateTokenSuccess tests successful token validation
func (suite *JWTServiceTestSuite) TestValidateTokenSuccess() {
	token, _ := suite.service.GenerateUserToken(suite.user)

	claims, err := suite.service.ValidateToken(token)

	suite.Empty( err.ErrCode)
	suite.Equal( suite.user.ID, claims["userId"])
	suite.Equal( suite.user.Username, claims["username"])
	suite.Equal( suite.user.Role, claims["role"])
}

// TestValidateTokenFailure tests validation failure for an invalid token
func (suite *JWTServiceTestSuite) TestValidateTokenFailure() {
	invalidToken := "invalid.token.string"

	_, err := suite.service.ValidateToken(invalidToken)

	suite.Equal(http.StatusUnauthorized, err.ErrCode)
	suite.Equal("Invalid token", err.ErrMessage)
}

// TestExpiredToken tests validation failure for an expired token
func (suite *JWTServiceTestSuite) TestExpiredToken() {
	// Manually create an expired token
	expiredClaims := &domain.Claims{
		UserId:   suite.user.ID,
		Username: suite.user.Username,
		Role:     suite.user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, _ := token.SignedString([]byte(suite.secret))

	_, err := suite.service.ValidateToken(tokenString)

	suite.Equal(http.StatusUnauthorized, err.ErrCode)
	suite.Equal("Invalid token", err.ErrMessage)
}

func TestJWTServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
}
