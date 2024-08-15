package infrastructure_test

import (
	"net/http"
	"task_managment_api/domain"
	"task_managment_api/infrastructure"
	"testing"
	"github.com/stretchr/testify/suite"
)

type PasswordServiceTestSuite struct {
	suite.Suite
	service infrastructure.PasswordService
}

func (suite *PasswordServiceTestSuite) SetupTest() {
	suite.service = infrastructure.NewPasswordService()
}

// TestHashPasswordSuccess tests successful password hashing
func (suite *PasswordServiceTestSuite) TestHashPasswordSuccess() {
	password := "securepassword123"
	hashedPassword, err := suite.service.HashPassword(password)

	suite.Empty( err.ErrCode)
	suite.NotEmpty( hashedPassword)
}


// TestVerifyPasswordSuccess tests successful password verification
func (suite *PasswordServiceTestSuite) TestVerifyPasswordSuccess() {
	password := "securepassword123"
	hashedPassword, err := suite.service.HashPassword(password)
	suite.Empty( err.ErrCode)

	user := domain.User{Password: hashedPassword}
	verificationErr := suite.service.VerifyPassword(user, password)

	suite.Empty(verificationErr.ErrCode)
}

// TestVerifyPasswordFailure tests password verification failure
func (suite *PasswordServiceTestSuite) TestVerifyPasswordFailure() {
	hashedPassword, err := suite.service.HashPassword("securepassword123")
	suite.Empty( err.ErrCode)

	user := domain.User{Password: hashedPassword}
	incorrectPassword := "wrongpassword"
	verificationErr := suite.service.VerifyPassword(user, incorrectPassword)

	suite.Equal( http.StatusUnauthorized, verificationErr.ErrCode)
	suite.Equal( "Invalid username or password", verificationErr.ErrMessage)
}

func TestPasswordServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordServiceTestSuite))
}
