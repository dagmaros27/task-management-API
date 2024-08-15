package domain

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// DomainTestSuite struct
type DomainTestSuite struct {
	suite.Suite
}

// SetupTest is run before each test
func (suite *DomainTestSuite) SetupTest() {
	// Setup code here if needed
}

// TestTaskInstantiation tests the instantiation of the Task model
func (suite *DomainTestSuite) TestTaskInstantiation() {
	task := Task{
		ID:          "123",
		Title:       "Test Task",
		Description: "This is a test task.",
		DueDate:     "2024-12-31",
		Status:      "Pending",
	}

	assert.NotNil(suite.T(), task)
	assert.Equal(suite.T(), "123", task.ID)
	assert.Equal(suite.T(), "Test Task", task.Title)
	assert.Equal(suite.T(), "This is a test task.", task.Description)
	assert.Equal(suite.T(), "2024-12-31", task.DueDate)
	assert.Equal(suite.T(), "Pending", task.Status)
}

// TestUserInstantiation tests the instantiation of the User model
func (suite *DomainTestSuite) TestUserInstantiation() {
	user := User{
		ID:       "456",
		Username: "testuser",
		Password: "securepassword",
		Role:     "user",
	}

	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "456", user.ID)
	assert.Equal(suite.T(), "testuser", user.Username)
	assert.Equal(suite.T(), "securepassword", user.Password)
	assert.Equal(suite.T(), "user", user.Role)
}

// TestUserToPromoteInstantiation tests the instantiation of the UserToPromote model
func (suite *DomainTestSuite) TestUserToPromoteInstantiation() {
	userToPromote := UserToPromote{
		Username: "testuser",
	}

	assert.NotNil(suite.T(), userToPromote)
	assert.Equal(suite.T(), "testuser", userToPromote.Username)
}

// TestClaimsInstantiation tests the instantiation of the Claims model
func (suite *DomainTestSuite) TestClaimsInstantiation() {
	claims := Claims{
		UserId: "789",
		Username: "testuser",
		Role: "user",
	}

	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), "789", claims.UserId)
	assert.Equal(suite.T(), "testuser", claims.Username)
	assert.Equal(suite.T(), "user", claims.Role)
}

// TestCustomErrorInstantiation tests the instantiation of the CustomError model
func (suite *DomainTestSuite) TestCustomErrorInstantiation() {
	customError := CustomError{
		ErrCode:    400,
		ErrMessage: "Bad Request",
	}

	assert.NotNil(suite.T(), customError)
	assert.Equal(suite.T(), 400, customError.ErrCode)
	assert.Equal(suite.T(), "Bad Request", customError.ErrMessage)
}

// Run the test suite
func TestDomainTestSuite(t *testing.T) {
	suite.Run(t, new(DomainTestSuite))
}
